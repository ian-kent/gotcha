package http

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	nethttp "net/http"
	neturl "net/url"
	"strings"

	"github.com/ian-kent/go-log/log"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	session    *Session
	writer     nethttp.ResponseWriter
	buffer     *bytes.Buffer
	headerSent bool

	Gzipped  bool
	gzwriter *gzip.Writer

	IsChunked bool

	IsEventStream bool
	esBufrw       *bufio.ReadWriter
	esConn        net.Conn

	Status  int
	Headers Headers
	Cookies Cookies
}

type Headers map[string][]string

func (h Headers) Add(name string, value string) {
	// TODO support same header multiple times
	h[name] = []string{value}
}

func (h Headers) Set(name string, value string) {
	h[name] = []string{value}
}

func (h Headers) Remove(name string) {
	delete(h, name)
}

type Cookies map[string]*nethttp.Cookie

func (c Cookies) Set(cookie *nethttp.Cookie) {
	c[cookie.Name] = cookie
}

func CreateResponse(session *Session, writer nethttp.ResponseWriter) *Response {
	return &Response{
		session:    session,
		writer:     writer,
		buffer:     &bytes.Buffer{},
		headerSent: false,

		Gzipped:       false,
		IsChunked:     false,
		IsEventStream: false,
		Status:        200,
		Headers:       make(Headers),
		Cookies:       make(Cookies),
	}
}

func (r *Response) GetWriter() nethttp.ResponseWriter {
	return r.writer
}

func toBase64(b []byte) string {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(b)
	encoder.Close()
	return buf.String()
}

func (r *Response) Gzip() {
	r.Gzipped = true
	r.Headers.Add("Content-Encoding", "gzip")
	r.gzwriter = gzip.NewWriter(r.writer)
}

func (r *Response) SessionID() string {
	if len(r.session.SessionData) > 0 && len(r.session.SessionID) == 0 {
		r.createSessionId()
	}
	return r.session.SessionID
}

func (r *Response) createSessionId() {
	bytes := make([]byte, 256)
	rand.Read(bytes)
	s, _ := bcrypt.GenerateFromPassword(bytes, 11)
	r.session.SessionID = string(s)
	log.Info("Generated session ID (__SID): %s", r.session.SessionID)
	r.Cookies.Set(&nethttp.Cookie{
		Name:  "__SID",
		Value: r.session.SessionID,
		Path:  "/",
	})
}

func (r *Response) writeSessionData() {
	b, err := json.Marshal(r.session.SessionData)

	if err != nil {
		r.session.RenderException(500, err)
		return
	}

	r.Cookies.Set(&nethttp.Cookie{
		Name:  "__SD",
		Value: toBase64(b),
		Path:  "/",
	})
}

func (r *Response) Write(bytes []byte) (int, error) {
	if r.headerSent {
		if r.Gzipped {
			return r.gzwriter.Write(bytes)
		} else {
			return r.writer.Write(bytes)
		}
	} else {
		return r.buffer.Write(bytes)
	}
}

func (r *Response) WriteText(text string) {
	r.Write([]byte(text))
}

func (r *Response) EventStream() chan []byte {
	c := make(chan []byte)

	r.IsEventStream = true
	r.Headers.Add("Content-Type", "text/event-stream")
	r.Headers.Add("Cache-Control", "no-cache")
	r.Headers.Add("Connection", "keep-alive")
	//r.Write([]byte("\n\n"))
	r.Send()

	hj, ok := r.writer.(nethttp.Hijacker)
	if !ok {
		log.Warn("Connection unsuitable for hijack")
		return nil
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		log.Warn("Connection hijack failed")
		return nil
	}

	r.esBufrw = bufrw
	r.esConn = conn

	go func() {
		for b := range c {
			if len(b) == 0 {
				log.Trace("Event stream ended")
				r.esConn.Close()
				break
			}

			lines := strings.Split(string(b), "\n")
			data := ""
			for _, l := range lines {
				data += "data: " + l + "\n"
			}
			data += "\n"

			sz := len(data) + 1
			log.Info("Event stream message is %d bytes", sz)
			size := fmt.Sprintf("%X", sz)
			r.esBufrw.Write([]byte(size + "\r\n"))

			lines = strings.Split(data, "\n")
			for _, ln := range lines {
				r.esBufrw.Write([]byte(ln + "\n"))
			}
			_, err := r.esBufrw.Write([]byte("\r\n"))

			if err != nil {
				log.Error("Error writing to connection: %s\n", err)
				r.esConn.Close()
				break
			}

			err = r.esBufrw.Flush()
			if err != nil {
				log.Error("Error flushing buffer: %s\n", err)
				r.esConn.Close()
				break
			}
		}
	}()
	return c
}

func (r *Response) Chunked() chan []byte {
	c := make(chan []byte)
	r.IsChunked = true
	r.Send()
	go func() {
		for b := range c {
			if len(b) == 0 {
				log.Trace("Chunk stream ended")
				if r.Gzipped {
					r.gzwriter.Close()
				}
				break
			}
			log.Trace("Writing chunk: %d bytes", len(b))
			if r.Gzipped {
				r.gzwriter.Write(b)
			} else {
				r.Write(b)
			}

			if f, ok := r.writer.(nethttp.Flusher); ok {
				f.Flush()
			}
		}
	}()
	return c
}

func (r *Response) Redirect(url *neturl.URL, status int) {
	r.Headers.Set("Location", url.String())
	r.Status = status
}

func (r *Response) Close() {
	if !r.IsChunked && r.Gzipped {
		r.gzwriter.Close()
	}
}

func (r *Response) Send() {
	if r.headerSent {
		return
	}
	r.headerSent = true

	r.SessionID()
	r.writeSessionData()

	for k, v := range r.Headers {
		for _, h := range v {
			log.Trace("Adding header [%s]: [%s]", k, h)
			r.writer.Header().Add(k, h)
		}
	}
	for _, c := range r.Cookies {
		nethttp.SetCookie(r.writer, c)
	}

	r.writer.WriteHeader(r.Status)

	if r.Gzipped {
		r.gzwriter.Write(r.buffer.Bytes())
	} else {
		r.writer.Write(r.buffer.Bytes())
	}

	if !r.IsChunked && !r.IsEventStream {
		if r.Gzipped {
			r.gzwriter.Close()
		}
	}
}
