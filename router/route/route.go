package route

import (
	"regexp"
)

type Route struct {
	Methods map[string]int
	Path    string
	Pattern *regexp.Regexp
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