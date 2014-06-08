package http

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/ian-kent/go-log/log"
	"github.com/ian-kent/gotcha/config"
	"github.com/ian-kent/gotcha/router/route"
	"html/template"
	nethttp "net/http"
	neturl "net/url"
	"runtime"
)

type Session struct {
	Config   *Config.Config
	Route    *route.Route
	Request  *Request
	Response *Response
	Stash    map[string]interface{}
	SessionID string
	SessionData map[string]string
}

func CreateSession(conf *Config.Config, request *nethttp.Request, writer nethttp.ResponseWriter) *Session {
	session := &Session{
		Config: conf,
		Stash:  make(map[string]interface{}, 0),
		SessionData: make(map[string]string),
	}

	// Somewhere to store internal stuff
	session.Stash["Gotcha"] = make(map[string]interface{}, 0)

	session.Request = CreateRequest(session, request)
	session.Response = CreateResponse(session, writer)

	session.loadSessionData()

	return session
}

func (s *Session) loadSessionData() {	
	if sid_cookie, ok := s.Request.Cookies["__SID"]; ok {
		s.SessionID = sid_cookie.Value
		log.Info("Retrieved session ID (__SID): %s", s.SessionID)
	}

	if sd_cookie, ok := s.Request.Cookies["__SD"]; ok {
		// FIXME deserialise
		s.SessionData["TEMP"] = sd_cookie.Value
		log.Info("Retrieved session data: %s", s.SessionData)
	}
}

func (session *Session) render(asset string) error {
	asset = "assets/templates/" + asset

	var t *template.Template

	c, ok := session.Config.Cache["template:"+asset]
	if !ok {
		log.Trace("Loading asset: %s", asset)
		a, err := session.Config.AssetLoader(asset)
		log.Trace("Creating template: %s", asset)
		t = template.New(asset)
		if err != nil || a == nil {
			log.Error("Failed loading template %s: %s", asset, err)
			return err
		}
		log.Trace("Parsing template: %s", asset)
		_, err = t.Parse(string(a))
		if err != nil {
			log.Error("Failed parsing template %s: %s", asset, err)
			return err
		}
		log.Trace("Template parsed successfully: %s", asset)
		session.Config.Cache["template:"+asset] = t
	} else {
		t = c.(*template.Template)
		log.Trace("Template loaded from cache: %s", asset)
	}

	var b bytes.Buffer
	err := t.Execute(&b, session.Stash)
	if err != nil {
		log.Error("Failed executing template %s: %s", asset, err)
		return err
	}

	_, err = session.Response.Write(b.Bytes())

	if err != nil {
		log.Error("Error writing output for template %s: %s", asset, err)
		return err
	}

	return nil
}

func (session *Session) RenderTemplate(asset string) (string, error) {
	asset = "assets/templates/" + asset

	var t *template.Template

	c, ok := session.Config.Cache["template:"+asset]
	if !ok {
		log.Trace("Loading asset: %s", asset)
		a, err := session.Config.AssetLoader(asset)
		log.Trace("Creating template: %s", asset)
		t = template.New(asset)
		if err != nil || a == nil {
			log.Error("Failed loading template %s: %s", asset, err)
			return "", err
		}
		log.Trace("Parsing template: %s", asset)
		_, err = t.Parse(string(a))
		if err != nil {
			log.Error("Failed parsing template %s: %s", asset, err)
			return "", err
		}
		log.Trace("Template parsed successfully: %s", asset)
		session.Config.Cache["template:"+asset] = t
	} else {
		t = c.(*template.Template)
		log.Trace("Template loaded from cache: %s", asset)
	}

	var b bytes.Buffer
	err := t.Execute(&b, session.Stash)
	if err != nil {
		log.Error("Failed executing template %s: %s", asset, err)
		return "", err
	}

	return b.String(), nil
}

func (session *Session) Render(asset string) {
	err := session.render(asset)
	if err != nil {
		session.RenderException(500, err)
	}
}

func (session *Session) RenderNotFound() {
	session.Stash["Gotcha"].(map[string]interface{})["Error"] = "Not found"
	session.Response.Status = 404

	e := session.render("notfound.html")
	if e != nil {
		log.Error("Error rendering not found page: %s", e)
	}
}

func (session *Session) RenderException(status int, err error) {
	key := uuid.NewUUID().String()
	session.Response.Status = status
	session.Stash["Gotcha"].(map[string]interface{})["Error"] = err.Error()
	session.Stash["Gotcha"].(map[string]interface{})["ErrorId"] = key

	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, true)
	session.Stash["Gotcha"].(map[string]interface{})["Stack"] = fmt.Sprintf("%s", buf[:n])

	log.Error("Exception %s: %s\n%s", key, err.Error(), session.Stash["Gotcha"].(map[string]interface{})["Stack"])

	e := session.render("error.html")
	if e != nil {
		log.Error("Error rendering error page: %s", e)
		session.Response.Write([]byte("[" + key + "] Internal Server Error: " + err.Error() + "\n"))
	}
}

func (session *Session) Redirect(url *neturl.URL) {
	log.Trace("Redirect to: %s", url)
	session.Response.Redirect(url, nethttp.StatusFound)
}
