package wrappers

import (
	"github.com/ian-kent/gotcha/http"
	"github.com/ian-kent/gotcha/router"
)

// An action to wrap other actions
func Foo(session *http.Session, f Router.HandlerFunc) {
	session.Stash["foo"] = "bar"
	// Call the nested action
	f(session)
}
