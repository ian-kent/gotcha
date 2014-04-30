package http

import (
	nethttp "net/http"
	neturl "net/url"
	"errors"
)

type Response struct {
	session *Session
	writer  nethttp.ResponseWriter
	code    int
}

func CreateResponse(session *Session, writer nethttp.ResponseWriter) *Response {
	return &Response{
		session: session,
		writer: writer,
	}
}

func (r *Response) Unwrap() nethttp.ResponseWriter {
	return r.writer
}

func (r *Response) NotFound() {
	r.code = 404
	r.session.RenderException(404, errors.New("Page not found"))
}

func (r *Response) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r *Response) WriteText(text string) {
	r.writer.Write([]byte(text))
}

func (r *Response) Status(status int) {
	r.code = status
	r.writer.WriteHeader(status)
}

func (r *Response) Headers() nethttp.Header {
	return r.writer.Header()
}

func (r *Response) Redirect(url *neturl.URL, status int) {
	r.code = status
	nethttp.Redirect(r.writer, r.session.Request.Unwrap(), url.String(), status)
}

func (r *Response) Code() int {
	return r.code
}
