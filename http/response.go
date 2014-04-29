package http

import (
	nethttp "net/http"
)

type Response struct {
	session *Session
	writer  nethttp.ResponseWriter
}

func CreateResponse(session *Session, writer nethttp.ResponseWriter) *Response {
	return &Response{
		session,
		writer,
	}
}

func (r *Response) NotFound() {
	nethttp.NotFound(r.writer, r.session.Request.Unwrap())
}

func (r *Response) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r *Response) WriteText(text string) {
	r.writer.Write([]byte(text))
}

func (r *Response) Status(status int) {
	r.writer.WriteHeader(status)
}

func (r *Response) Headers() nethttp.Header {
	return r.writer.Header()
}
