package http

import (
	nethttp "net/http"
	neturl "net/url"
	"bytes"
)

type Response struct {
	session *Session
	writer  nethttp.ResponseWriter
	buffer  *bytes.Buffer

	Status    int
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
		session: session,
		writer: writer,
		buffer: &bytes.Buffer{},

		Status: 200,
		Headers: make(Headers),
		Cookies: make(Cookies),
	}
}

func (r *Response) Write(bytes []byte) (int, error) {
	return r.buffer.Write(bytes)
}

func (r *Response) WriteText(text string) {
	r.buffer.Write([]byte(text))
}

func (r *Response) Redirect(url *neturl.URL, status int) {
	r.Headers.Set("Location", url.String())
	r.Status = status
}

func (r *Response) Send() {
	for k, v := range r.Headers {
		for _, h := range v {
			// log.Printf("[TRACE] Adding header [%s]: [%s]", k, h)
			r.writer.Header().Add(k, h)
		}
	}
	for _, c := range r.Cookies {
		nethttp.SetCookie(r.writer, c)
	}
	r.writer.WriteHeader(r.Status)
	r.writer.Write(r.buffer.Bytes())
}
