package http

import (
	"bytes"
	"github.com/ian-kent/gotcha/config"
	"html/template"
	"log"
	nethttp "net/http"
	neturl "net/url"
)

type Session struct {
	Config   *Config.Config
	Request  *Request
	Response *Response
	Stash    map[string]interface{}
}

func CreateSession(conf *Config.Config, request *nethttp.Request, writer nethttp.ResponseWriter) *Session {
	session := &Session{
		Config: conf,
		Stash:  make(map[string]interface{}, 0),
	}

	session.Request = CreateRequest(session, request)
	session.Response = CreateResponse(session, writer)

	return session
}

func (session *Session) render(asset string) error {
	asset = "assets/templates/" + asset

	var t *template.Template

	c, ok := session.Config.Cache["template:"+asset]
	if !ok {
		log.Printf("Parsing template: %s", asset)
		t = template.New(asset)
		a, err := session.Config.AssetLoader(asset)
		if err != nil {
			log.Printf("Failed loading template %s: %s", asset, err)
			return err
		}
		_, err = t.Parse(string(a))
		if err != nil {
			log.Printf("Failed parsing template %s: %s", asset, err)
			return err
		}
		log.Printf("Template parsed successfully: %s", asset)
		session.Config.Cache["template:"+asset] = t
	} else {
		t = c.(*template.Template)
		log.Printf("Template loaded from cache: %s", asset)
	}

	var b bytes.Buffer
	err := t.Execute(&b, session.Stash)
	if err != nil {
		log.Printf("Failed executing template %s: %s", asset, err)
		return err
	}

	_, err = session.Response.Write(b.Bytes())

	if err != nil {
		log.Printf("Error writing output for template %s: %s", asset, err)
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

func (session *Session) Redirect(url *neturl.URL) {
	session.Response.Redirect(url, nethttp.StatusFound)
}
