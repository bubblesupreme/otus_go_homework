package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"

	log "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/gorilla/mux"
)

type Server struct {
	port   int
	host   string
	router *mux.Router
	server *http.Server
	app    Application
}

type Application interface{}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Home"))
	if err != nil {
		log.Error(err.Error(), nil)
	}
}

func helloHandler(w http.ResponseWriter, h *http.Request) {
	_, err := w.Write([]byte("Hello world!"))
	if err != nil {
		log.Error(err.Error(), nil)
	}
}

func NewServer(app Application, port int, host string) *Server {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/hello", helloHandler)
	r.Use(loggingMiddleware)
	http.Handle("/", r)
	return &Server{
		port:   port,
		host:   host,
		router: r,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{
		Addr:    net.JoinHostPort(s.host, strconv.Itoa(s.port)),
		Handler: s.router,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
