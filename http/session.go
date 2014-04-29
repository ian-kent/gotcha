package http

import (
	nethttp "net/http"
	"github.com/ian-kent/Go-Gotcha/config"
)

type Session struct {
	Config *Config.Config
	Request *Request
	Response *Response
}

func CreateSession(conf *Config.Config, request *nethttp.Request, writer nethttp.ResponseWriter) *Session {
	session := &Session{
		Config: conf,
	}

	session.Request = CreateRequest(session, request)
	session.Response = CreateResponse(session, writer)

	return session
}
