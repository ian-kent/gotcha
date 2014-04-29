package Router

import (
	"github.com/ian-kent/gotcha/config"
	"github.com/ian-kent/gotcha/http"
	"github.com/ian-kent/gotcha/mime"
	"log"
	nethttp "net/http"
	"regexp"
)

// http://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc

type Route struct {
	Methods map[string]int
	Path    string
	Pattern *regexp.Regexp
	Handler HandlerFunc
}

func (r1 *Route) Equals(r2 *Route) bool {
	if r1.Path != r2.Path {
		return false
	}
	for k, _ := range r1.Methods {
		_, ok := r2.Methods[k]
		if !ok {
			return false
		}
	}
	for k, _ := range r2.Methods {
		_, ok := r1.Methods[k]
		if !ok {
			return false
		}
	}
	return true
}

type HandlerFunc func(*http.Session, *Route)

func (f HandlerFunc) ServeHTTP(session *http.Session, route *Route) {
	f(session, route)
}

type Router struct {
	Config *Config.Config
	routes []*Route
}

func Create(config *Config.Config) *Router {
	return &Router{
		Config: config,
	}
}

func (h *Router) Routes() []*Route {
	return h.routes
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
	return func(session *http.Session, route *Route) {
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
			session.Response.NotFound()
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
	h.routes = append(h.routes, &Route{m, path, pattern, handler})
}

func (h *Router) HandleFunc(methods []string, path string, handler func(*http.Session, *Route)) {
	pattern := regexp.MustCompile("^" + path + "$")
	m := make(map[string]int, 0)
	for _, v := range methods {
		m[v] = 1
	}
	h.routes = append(h.routes, &Route{m, path, pattern, HandlerFunc(handler)})
}

func (h *Router) Serve(session *http.Session) {
	for _, route := range h.routes {
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
				route.Handler.ServeHTTP(session, route)
				return
			}
		}
	}
	// no pattern matched; send 404 response
	session.Response.NotFound()
}

func (h *Router) ServeHTTP(w nethttp.ResponseWriter, r *nethttp.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	session := http.CreateSession(h.Config, r, w)
	h.Serve(session)
}
