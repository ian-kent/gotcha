package http

import (
	"io"
	"mime/multipart"
	nethttp "net/http"
	"net/url"
)

type Request struct {
	Session    *Session
	req        *nethttp.Request
	URL        *url.URL
	Method     string
	Cookies    map[string]*nethttp.Cookie
	RemoteAddr string
}

func CreateRequest(session *Session, request *nethttp.Request) *Request {
	req := &Request{
		Session:    session,
		req:        request,
		URL:        request.URL,
		Method:     request.Method,
		Cookies:    make(map[string]*nethttp.Cookie),
		RemoteAddr: request.RemoteAddr,
	}

	for _, cookie := range request.Cookies() {
		req.Cookies[cookie.Name] = cookie
	}

	return req
}

func (r *Request) Form() url.Values {
	r.req.ParseForm()
	return r.req.Form
}

func (r *Request) PostForm() url.Values {
	r.req.ParseForm()
	return r.req.PostForm
}

func (r *Request) MultipartForm() *multipart.Form {
	r.req.ParseMultipartForm(1024000) // FIXME
	return r.req.MultipartForm
}

func (r *Request) File(input string) (multipart.File, *multipart.FileHeader, error) {
	return r.req.FormFile(input)
}

func (r *Request) Header() nethttp.Header {
	return r.req.Header
}

func (r *Request) Body() io.ReadCloser {
	return r.req.Body
}
