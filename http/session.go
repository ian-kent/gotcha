package http

import (
	nethttp "net/http"
	"html/template"
	"github.com/ian-kent/Go-Gotcha/config"
	"bytes"
)

type Session struct {
	Config *Config.Config
	Request *Request
	Response *Response
	Stash map[string]interface{}
}

func CreateSession(conf *Config.Config, request *nethttp.Request, writer nethttp.ResponseWriter) *Session {
	session := &Session{
		Config: conf,
	}

	session.Request = CreateRequest(session, request)
	session.Response = CreateResponse(session, writer)

	return session
}

func (session *Session) Render(asset string) error {
	asset = "assets/templates/" + asset

	t := template.New(asset)
	a, err := session.Config.AssetLoader(asset)
	if err != nil {
		return err
	}

	t.Parse(string(a))
	var b bytes.Buffer
	t.Execute(&b, session.Stash)
	_, err = session.Response.Write(b.Bytes())

	if err != nil {
		return err
	}

	return nil
}
