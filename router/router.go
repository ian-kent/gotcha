package Router

import (
	"github.com/ian-kent/Go-Gotcha/http"
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

type HandlerFunc func(*http.Response, *http.Request, *Route)

//type Handler http.Handler
func (f HandlerFunc) ServeHTTP(w *http.Response, r *http.Request, route *Route) {
	f(w, r, route)
}

type Router struct {
	routes []*Route
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

func Static(filename string) HandlerFunc {
	return func(w *http.Response, r *http.Request, route *Route) {
		w.NotFound()
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

func (h *Router) HandleFunc(methods []string, path string, handler func(*http.Response, *http.Request, *Route)) {
	pattern := regexp.MustCompile("^" + path + "$")
	m := make(map[string]int, 0)
	for _, v := range methods {
		m[v] = 1
	}
	h.routes = append(h.routes, &Route{m, path, pattern, HandlerFunc(handler)})
}

func (h *Router) Serve(w *http.Response, r *http.Request) {
	for _, route := range h.routes {
		if route.Pattern.MatchString(r.URL.Path) {
			_, ok := route.Methods[r.Method]
			if ok {
				route.Handler.ServeHTTP(w, r, route)
				return
			}
		}
	}
	// no pattern matched; send 404 response
	w.NotFound()
}

func (h *Router) ServeHTTP(w nethttp.ResponseWriter, r *nethttp.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	req := http.CreateRequest(r)
	res := http.CreateResponse(req, w)
	h.Serve(res, req)
}
