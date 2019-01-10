/*
 * File share api
 *
 * File share api.
 *
 * API version: 2.0.0
 * Contact: enponsba@gmail.com
 */

package raspi

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{

	Route{
		"GetFileById",
		strings.ToUpper("Get"),
		"/api/v1/file/{fileId}",
		GetFileById,
	},

	Route{
		"ListFiles",
		strings.ToUpper("Get"),
		"/api/v1/file",
		ListFiles,
	},

	Route{
		"NewFile",
		strings.ToUpper("Post"),
		"/api/v1/file",
		NewFile,
	},
}
