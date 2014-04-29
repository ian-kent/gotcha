package http

import (
	nethttp "net/http"
	neturl "net/url"
)

type Request struct {
	Session *Session
	req     *nethttp.Request
	URL     *neturl.URL
	Method  string
}

func CreateRequest(session *Session, request *nethttp.Request) *Request {
	return &Request{
		session,
		request,
		request.URL,
		request.Method,
	}
}

// FIXME
func (r *Request) Unwrap() *nethttp.Request {
	return r.req
}
