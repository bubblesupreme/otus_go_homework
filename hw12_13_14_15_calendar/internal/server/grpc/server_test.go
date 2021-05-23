package internalgrpc

import (
	"context"
	"net"
	"strconv"
	"testing"

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

func doTest(t *testing.T, testFn func(t *testing.T, server server.Server, client server.TestClient)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.Dial(net.JoinHostPort(host, strconv.Itoa(port)), grpc.WithInsecure())
	assert.NoError(t, err)
	client := TestClientGRPC{
		grpcClient: eventspb.NewEventServiceClient(conn),
	}

	storageMem, err := server.GetMemoryStorage()
	assert.NoError(t, err)
	aMem := app.NewApp(storageMem)
	sMem, err := NewServer(ctx, aMem, port, host)
	assert.NoError(t, err)
	testFn(t, sMem, client)

	// TODO: docker compose
	// storageSQL, err := server.GetSQLStorage()
	// assert.NoError(t, err)
	// aSQL := app.NewApp(storageSQL)
	// sSQL, err := NewServer(ctx, aSQL, port, host)
	// assert.NoError(t, err)
	// testFn(t, sSQL, client)
}

func TestEmptyStorage(t *testing.T) {
	doTest(t, server.EmptyStorageImpl)
}

func TestCommon(t *testing.T) {
	doTest(t, server.CommonImpl)
}
