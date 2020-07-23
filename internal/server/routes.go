package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RouteHandler type.
type RouteHandler func(*server) http.Handler

// Route struct.
type Route struct {
	Method  string
	Path    string
	Handler RouteHandler
}

var (
	// RouteTable is route definition.
	RouteTable = []Route{
		{
			Method: "GET", Path: "/_healthz", Handler: healthzHandler,
		},
	}

	SocksRouteTable = []Route{
		{
			Method: "GET", Path: "/socks_healthz", Handler: socksHealthzHandler,
		},
	}
)

func healthzHandler(s *server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(map[string]string{"name": "wireguard", "status": "OK"})
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(b))
	})
}

func socksHealthzHandler(s *server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(map[string]string{"name": "socks proxy", "status": "OK"})
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(b))
	})
}
