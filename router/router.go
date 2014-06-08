package Router

import (
	"errors"
	"github.com/ian-kent/go-log/log"
	"github.com/ian-kent/gotcha/config"
	"github.com/ian-kent/gotcha/events"
	"github.com/ian-kent/gotcha/http"
	"github.com/ian-kent/gotcha/mime"
	"github.com/ian-kent/gotcha/router/route"
	nethttp "net/http"
	"regexp"
	"time"
)

// http://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc

type HandlerFunc func(*http.Session)

func (f HandlerFunc) ServeHTTP(session *http.Session) {
	f(session)
}

type Router struct {
	Config *Config.Config
	Routes []*Route
}

type Route struct {
	Route *route.Route
	Handler HandlerFunc
}

func Create(config *Config.Config) *Router {
	return &Router{
		Config: config,
		Routes: make([]*Route,0),
	}
}

func (h *Router) Get(pattern string, handler HandlerFunc) *Router {
	h.Handler([]string{"GET"}, pattern, handler)
	return h
}

func (h *Router) Post(pattern string, handler HandlerFunc) *Router {
	h.Handler([]string{"POST"}, pattern, handler)
	return h
}

func (h *Router) Put(pattern string, handler HandlerFunc) *Router {
	h.Handler([]string{"PUT"}, pattern, handler)
	return h
}

func (h *Router) Delete(pattern string, handler HandlerFunc) *Router {
	h.Handler([]string{"DELETE"}, pattern, handler)
	return h
}

func (h *Router) Patch(pattern string, handler HandlerFunc) *Router {
	h.Handler([]string{"PATCH"}, pattern, handler)
	return h
}

func (h *Router) Options(pattern string, handler HandlerFunc) *Router {
	h.Handler([]string{"OPTIONS"}, pattern, handler)
	return h
}

func (h *Router) Static(filename string) HandlerFunc {
	return func(session *http.Session) {
		// TODO beware of ..?
		re := regexp.MustCompile("{{(\\w+)}}")
		fcopy := re.ReplaceAllStringFunc(filename, func(m string) string {
			parts := re.FindStringSubmatch(m)
			log.Trace("Found var: %s; name: %s", m, parts[1])
			if val, ok := session.Stash[parts[1]]; ok {
				log.Trace("Value found in stash for %s: %s", parts[1], val)
				return val.(string)
			}
			log.Trace("No value found in stash for var: %s", parts[1])
			return m
		})
		asset, err := h.Config.AssetLoader(fcopy)
		if err != nil {
			log.Debug("Static file not found: %s", fcopy)
			session.RenderNotFound()
		} else {
			m := MIME.TypeFromFilename(fcopy)
			if len(m) > 0 {
				log.Debug("Setting Content-Type: %s", m)
				session.Response.Headers.Add("Content-Type", m[0])
			}
			session.Response.Write(asset)
		}
	}
}

func PatternToRegex(pattern string) *regexp.Regexp {
	return regexp.MustCompile("^" + pattern + "$")
}

func (h *Router) Handler(methods []string, path string, handler HandlerFunc) *Router {
	pattern := PatternToRegex(path)
	m := make(map[string]int, 0)
	for _, v := range methods {
		m[v] = 1
	}
	h.Routes = append(h.Routes, &Route{Route:&route.Route{m, path, pattern},Handler:handler})
	return h
}

func (h *Router) HandleFunc(methods []string, path string, handler func(*http.Session)) *Router {
	pattern := PatternToRegex(path)
	m := make(map[string]int, 0)
	for _, v := range methods {
		m[v] = 1
	}
	h.Routes = append(h.Routes, &Route{Route:&route.Route{m, path, pattern},Handler:HandlerFunc(handler)})
	return h
}

func (h *Router) Serve(session *http.Session) {
	for _, route := range h.Routes {
		if matches := route.Route.Pattern.FindStringSubmatch(session.Request.URL.Path); len(matches) > 0 {
			_, ok := route.Route.Methods[session.Request.Method]
			if ok {
				for i, named := range route.Route.Pattern.SubexpNames() {
					if len(named) > 0 {
						log.Trace("Matched named pattern '%s': %s", named, matches[i])
						session.Stash[named] = matches[i]
					}
				}
				session.Route = route.Route
				defer func() {
					if e := recover(); e != nil {
						switch e.(type) {
						case string:
							session.RenderException(500, errors.New(e.(string)))
						default:
							session.RenderException(500, e.(error))
						}
						session.Response.Send()
					}
				}()
				// func() will be executed only if *all* event handlers call next()
				h.Config.Events.Emit(session, events.BeforeHandler, func() {
					route.Handler.ServeHTTP(session)
					h.Config.Events.Emit(session, events.AfterHandler, func() {
						session.Response.Send()
					})
				})
				return
			}
		}
	}

	// no pattern matched; send 404 response	
	h.Config.Events.Emit(session, events.BeforeHandler, func() {
		session.RenderNotFound()
		h.Config.Events.Emit(session, events.AfterHandler, func() {
			session.Response.Send()
		})
	})
}

func (h *Router) ServeHTTP(w nethttp.ResponseWriter, r *nethttp.Request) {
	session := http.CreateSession(h.Config, r, w)
	tStart := time.Now().UnixNano()

	h.Serve(session)

	t := float64(time.Now().UnixNano()-tStart) / 100000 // ms
	log.Printf("%s %s (%3.2fms) (%d)", r.Method, r.URL, t, session.Response.Status)
	h.Config.Events.Emit(session, events.AfterResponse, func() {})
}
