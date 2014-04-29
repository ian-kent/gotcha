package http

import (
	nethttp "net/http"
	"html/template"
	"github.com/ian-kent/gotcha/config"
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
		Stash: make(map[string]interface{},0),
	}

	session.Request = CreateRequest(session, request)
	session.Response = CreateResponse(session, writer)

	return session
}

func (session *Session) render(asset string) error {
	asset = "assets/templates/" + asset

	t := template.New(asset)
	a, err := session.Config.AssetLoader(asset)
	if err != nil {
		return err
	}

	t.Parse(string(a))
	var b bytes.Buffer
	err = t.Execute(&b, session.Stash)
	if err != nil {
		return err
	}

	_, err = session.Response.Write(b.Bytes())

	if err != nil {
		return err
	}

	return nil
}

func (session *Session) Render(asset string) {
	err := session.render(asset)
	if err != nil {
		session.RenderException(500, err)
	}
}

func (session *Session) RenderException(status int, err error) {
	session.Response.Status(status)
	session.Response.Write([]byte("Internal Server Error: "))
	session.Response.Write([]byte(err.Error()))
}
