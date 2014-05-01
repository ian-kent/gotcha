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
	Cookies map[string]*nethttp.Cookie
}

func CreateRequest(session *Session, request *nethttp.Request) *Request {
	req := &Request{
		Session: session,
		req: request,
		URL: request.URL,
		Method: request.Method,
		Cookies: make(map[string]*nethttp.Cookie),
	}

	for _, cookie := range request.Cookies() {
		req.Cookies[cookie.Name] = cookie
	}

	return req
}
