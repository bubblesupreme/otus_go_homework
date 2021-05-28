package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"

	eventspb "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	server *http.Server
	mux    *runtime.ServeMux
}

func NewServer(ctx context.Context, app *app.App, port int, host string) (*Server, error) {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/json", &runtime.JSONBuiltin{}))
	err := eventspb.RegisterEventServiceHandlerServer(ctx, mux, app)
	s := http.Server{
		Addr:    net.JoinHostPort("", strconv.Itoa(port)),
		Handler: loggingMiddleware(mux),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	return &Server{
		server: &s,
		mux:    mux,
	}, err
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.server.Close()
}
