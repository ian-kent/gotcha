package Router

import (
	"github.com/ian-kent/gotcha/config"
	"github.com/ian-kent/gotcha/http"
	"github.com/ian-kent/gotcha/mime"
	"github.com/ian-kent/gotcha/events"
	"github.com/ian-kent/gotcha/router/route"
	"log"
	nethttp "net/http"
	"regexp"
	"time"
	"errors"
)

// http://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc

type HandlerFunc func(*http.Session)

func (f HandlerFunc) ServeHTTP(session *http.Session) {
	f(session)
}

type Router struct {
	Config *Config.Config
	Routes map[*route.Route]HandlerFunc
}

func Create(config *Config.Config) *Router {
	return &Router{
		Config: config,
		Routes: make(map[*route.Route]HandlerFunc),
	}
}

func (h *Router) Get(pattern string, handler HandlerFunc) {
	h.Handler([]string{"GET"}, pattern, handler)
}

func (h *Router) Post(pattern string, handler HandlerFunc) {
	h.Handler([]string{"POST"}, pattern, handler)
}

func (h *Router) Put(pattern string, handler HandlerFunc) {
	h.Handler([]string{"PUT"}, pattern, handler)
}

func (h *Router) Delete(pattern string, handler HandlerFunc) {
	h.Handler([]string{"DELETE"}, pattern, handler)
}

func (h *Router) Patch(pattern string, handler HandlerFunc) {
	h.Handler([]string{"PATCH"}, pattern, handler)
}

func (h *Router) Options(pattern string, handler HandlerFunc) {
	h.Handler([]string{"OPTIONS"}, pattern, handler)
}

func (h *Router) Static(filename string) HandlerFunc {
	return func(session *http.Session) {
		// TODO beware of ..?
		// TODO re-add log lines using log4go
		re := regexp.MustCompile("{{(\\w+)}}")
		fcopy := re.ReplaceAllStringFunc(filename, func(m string) string {
			//log.Printf("Found var: %s", m)
			parts := re.FindStringSubmatch(m)
			//log.Printf("Found var: %s; name: %s", m, parts[1])
			if val, ok := session.Stash[parts[1]]; ok {
				//log.Printf("Value found in stash for %s: %s", parts[1], val)
				return val.(string)
			}
			log.Printf("No value found in stash for var: %s", parts[1])
			return m
		})
		asset, err := h.Config.AssetLoader(fcopy)
		if err != nil {
			log.Printf("Static file not found: " + fcopy)
			session.RenderNotFound()
		} else {
			m := MIME.TypeFromFilename(fcopy)
			if len(m) > 0 {
				session.Response.Headers().Add("Content-Type", m[0])
			}
			session.Response.Write(asset)
		}
	}
}

func (h *Router) Handler(methods []string, path string, handler HandlerFunc) {
	pattern := regexp.MustCompile("^" + path + "$")
	m := make(map[string]int, 0)
	for _, v := range methods {
		m[v] = 1
	}
	h.Routes[&route.Route{m, path, pattern}] = handler
}

func (h *Router) HandleFunc(methods []string, path string, handler func(*http.Session)) {
	pattern := regexp.MustCompile("^" + path + "$")
	m := make(map[string]int, 0)
	for _, v := range methods {
		m[v] = 1
	}
	h.Routes[&route.Route{m, path, pattern}] = HandlerFunc(handler)
}

func (h *Router) Serve(session *http.Session) (t float64) {
	tStart := time.Now().UnixNano()
	for route, handler := range h.Routes {
		if matches := route.Pattern.FindStringSubmatch(session.Request.URL.Path); len(matches) > 0 {
			_, ok := route.Methods[session.Request.Method]
			if ok {
				for i, named := range route.Pattern.SubexpNames() {
					if len(named) > 0 {
						// TODO log4go
						//log.Printf("Matched named pattern '%s': %s", named, matches[i])
						session.Stash[named] = matches[i]
					}
				}
				session.Route = route
				defer func() {
					if e := recover(); e != nil {
						session.RenderException(500, errors.New(e.(string)))
						t = float64(time.Now().UnixNano()-tStart) / 100000
					}
				}()
				h.Config.Events <- &events.Event{events.AfterHandler}
				handler.ServeHTTP(session)
				return float64(time.Now().UnixNano()-tStart) / 100000
			}
		}
	}
	// no pattern matched; send 404 response
	session.RenderNotFound()
	return float64(time.Now().UnixNano()-tStart) / 100000
}

func (h *Router) ServeHTTP(w nethttp.ResponseWriter, r *nethttp.Request) {
	session := http.CreateSession(h.Config, r, w)
	ms := h.Serve(session)
	log.Printf("%s %s (%3.2fms) (%d)\n", r.Method, r.URL, ms, session.Response.Code())
}
