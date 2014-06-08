package http

import (
	"bytes"
	"github.com/ian-kent/go-log/log"
	nethttp "net/http"
	neturl "net/url"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"compress/gzip"
	"io"
)

type Response struct {
	session    *Session
	writer     nethttp.ResponseWriter
	buffer     *bytes.Buffer
	headerSent bool

	Gzipped bool
	gz *gzip.Writer
	IsChunked bool

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

		Gzipped: false,
		IsChunked: false,
		Status:  200,
		Headers: make(Headers),
		Cookies: make(Cookies),
	}
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
	})
}

func (r *Response) Write(bytes []byte) (int, error) {
	return r.buffer.Write(bytes)
}

func (r *Response) WriteText(text string) {
	r.buffer.Write([]byte(text))
}

type gzipResponseWriter struct {
	io.Writer
	nethttp.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (r *Response) Gzip() {
	r.Gzipped = true
	r.gz = gzip.NewWriter(r.writer)
	r.writer = gzipResponseWriter{Writer: r.gz, ResponseWriter: r.writer}
	r.Headers.Add("Content-Encoding", "gzip")
}

func (r *Response) Chunked() chan []byte {
	c := make(chan []byte)
	r.Send()
	r.IsChunked = true
	go func() {
		for b := range c {
			if len(b) == 0 {
				log.Trace("Chunk stream ended")
				if r.Gzipped {
					r.gz.Close()
				}
				break
			}
			log.Trace("Writing chunk: %d bytes", len(b))
			r.writer.Write(b)
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

func (r *Response) Send() {
	if r.headerSent {
		return
	}
	r.headerSent = true

	if len(r.session.SessionData) > 0 && len(r.session.SessionID) == 0 {
		r.createSessionId()
		//r.writeSessionData()
	}

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
	r.writer.Write(r.buffer.Bytes())

	if !r.IsChunked && r.Gzipped {
		r.gz.Close()
	}
}
