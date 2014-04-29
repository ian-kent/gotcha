package Router

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"regexp"
	"testing"
)

func TestRouter(t *testing.T) {
	router := &Router{}
	assert.NotNil(t, router)

	assert.Equal(t, len(router.Routes()), 0)

	f := func(w http.ResponseWriter, r *http.Request, route *Route) {}
	router.Get("/", f)
	assert.Equal(t, len(router.Routes()), 1)
	assert.Equal(t, true, router.Routes()[0].Equals(&Route{map[string]int{"GET": 1}, "/", regexp.MustCompile("/"), f}))

	f = func(w http.ResponseWriter, r *http.Request, route *Route) {}
	router.Post("/", f)
	assert.Equal(t, len(router.Routes()), 2)
	assert.Equal(t, true, router.Routes()[0].Equals(&Route{map[string]int{"GET": 1}, "/", regexp.MustCompile("/"), f}))
	assert.Equal(t, true, router.Routes()[1].Equals(&Route{map[string]int{"POST": 1}, "/", regexp.MustCompile("/"), f}))
}
