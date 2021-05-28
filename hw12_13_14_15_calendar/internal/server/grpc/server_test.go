package internalgrpc

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	eventspb "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/api"
	"google.golang.org/grpc"

	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/server"
	"github.com/stretchr/testify/assert"
)

const (
	host = "localhost"
	port = 8080
)

type TestClientGRPC struct {
	grpcClient eventspb.EventServiceClient
}

func (t TestClientGRPC) CreateEvent(ctx context.Context, in *eventspb.Event) (*eventspb.EventID, error) {
	return t.grpcClient.CreateEvent(ctx, in)
}

func (t TestClientGRPC) UpdateEvent(ctx context.Context, in *eventspb.Event) (*eventspb.Empty, error) {
	return t.grpcClient.UpdateEvent(ctx, in)
}

func (t TestClientGRPC) RemoveEvent(ctx context.Context, in *eventspb.EventID) (*eventspb.Empty, error) {
	return t.grpcClient.RemoveEvent(ctx, in)
}

func (t TestClientGRPC) GetDayEvents(ctx context.Context, in *eventspb.Time) (*eventspb.Events, error) {
	return t.grpcClient.GetDayEvents(ctx, in)
}

func (t TestClientGRPC) GetWeekEvents(ctx context.Context, in *eventspb.Time) (*eventspb.Events, error) {
	return t.grpcClient.GetWeekEvents(ctx, in)
}

func (t TestClientGRPC) GetMonthEvents(ctx context.Context, in *eventspb.Time) (*eventspb.Events, error) {
	return t.grpcClient.GetMonthEvents(ctx, in)
}

func TestGRPC(t *testing.T) {
	t.Run("EmptyMemoryStorage", func(t *testing.T) { doTestMemory(t, server.EmptyStorageImpl) })
	t.Run("EmptySQLStorage", func(t *testing.T) { /*doTestSQL(t, server.EmptyStorageImpl)*/ })
	t.Run("CommonMemoryStorage", func(t *testing.T) { doTestMemory(t, server.CommonImpl) })
	t.Run("CommonSQLStorage", func(t *testing.T) { /*doTestSQL(t, server.CommonImpl)*/ })
}

func doTestMemory(t *testing.T, testFn func(t *testing.T, client server.TestClient)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.Dial(net.JoinHostPort(host, strconv.Itoa(port)), grpc.WithInsecure())
	assert.NoError(t, err)
	client := TestClientGRPC{
		grpcClient: eventspb.NewEventServiceClient(conn),
	}

	storage, err := server.GetMemoryStorage()
	assert.NoError(t, err)
	a := app.NewApp(storage)
	srv, err := NewServer(ctx, a, port, host)
	assert.NoError(t, err)
	go func() {
		assert.NoError(t, srv.Start())
	}()
	defer func() {
		assert.NoError(t, srv.Stop())
	}()
	time.Sleep(5 * time.Second) // waiting for server start
	testFn(t, client)
}

func doTestSQL(t *testing.T, testFn func(t *testing.T, client server.TestClient)) { //nolint:deadcode,unused
	// TODO: docker compose
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.Dial(net.JoinHostPort(host, strconv.Itoa(port)), grpc.WithInsecure())
	assert.NoError(t, err)
	client := TestClientGRPC{
		grpcClient: eventspb.NewEventServiceClient(conn),
	}
	storage, err := server.GetSQLStorage()
	assert.NoError(t, err)
	a := app.NewApp(storage)
	srv, err := NewServer(ctx, a, port, host)
	assert.NoError(t, err)
	go func() {
		assert.NoError(t, srv.Start())
	}()
	defer func() {
		assert.NoError(t, srv.Stop())
	}()
	time.Sleep(5 * time.Second) // waiting for server start
	testFn(t, client)
}
