package common

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

func (r *Router) AddRoute(name, method, pattern string, handlerfunc http.HandlerFunc) {
	r.
		Methods(method).
		Path(pattern).
		Name(name).
		HandlerFunc(handlerfunc)
}

func NewRouter() *Router {
	router := mux.NewRouter().StrictSlash(true)
	return &Router{Router: router}
}
