package internalgrpc

import (
	"context"
	"fmt"
	"net"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"

	eventspb "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/api"
	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	log "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	port   int
	host   string
	app    *app.App
	server *grpc.Server
}

func NewServer(_ context.Context, app *app.App, port int, host string) (*Server, error) {
	logEntry := logrus.NewEntry(log.GetLogger())
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_logrus.StreamServerInterceptor(logEntry),
	)))
	reflection.Register(grpcServer)
	eventspb.RegisterEventServiceServer(grpcServer, app)

	return &Server{
		port:   port,
		host:   host,
		app:    app,
		server: grpcServer,
	}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		log.Error("failed to listen", log.Fields{
			"host": s.host,
			"port": s.port,
		})
		return err
	}
	return s.server.Serve(lis)
}

func (s *Server) Stop() error {
	s.server.Stop()
	return nil
}
