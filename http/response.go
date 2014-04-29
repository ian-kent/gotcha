package http

import(
	nethttp "net/http"
)

type Response struct {
	writer nethttp.ResponseWriter
	request *Request
}

func CreateResponse(request *Request, writer nethttp.ResponseWriter) *Response {
	return &Response{
		writer,
		request,
	}
}

func (r *Response) NotFound() {
	nethttp.NotFound(r.writer, r.request.Unwrap())
}

func (r *Response) Write(bytes []byte) {
	r.writer.Write(bytes)
}
