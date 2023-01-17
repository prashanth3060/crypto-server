package httpserver

import (
	"context"
	"net"
	"net/http"

	"github.com/gookit/slog"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	router *mux.Router
	server *http.Server
}

func NewHTTPServer(ctx context.Context, r ...Routes) *HTTPServer {
	router := mux.NewRouter()
	sr := router.PathPrefix("/api/v1/").Subrouter()
	server := &HTTPServer{
		router: sr,
		server: &http.Server{
			Addr:    net.JoinHostPort("", "8080"),
			Handler: router,
		},
	}
	server.addRoutes(r...)
	return server
}

func (s HTTPServer) addRoutes(r ...Routes) {
	var allRoutes Routes

	for _, route := range r {
		allRoutes = append(allRoutes, route...)
	}

	slog.Info("registering routes")
	for _, r := range allRoutes {
		slog.Infof("adding route, method: %s, path: %s", r.Verb, r.Path)
		s.router.Methods(r.Verb).Path(r.Path).HandlerFunc(r.Fn)
	}
}

func (s HTTPServer) Start() error {
	slog.Infof("starting http server, PORT: %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s HTTPServer) Close() error {
	return s.server.Close()
}
