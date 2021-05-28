package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	eventspb "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/server"
	"github.com/stretchr/testify/assert"
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

type TestClientHTTP struct {
	addr string
}

func (t TestClientHTTP) CreateEvent(_ context.Context, in *eventspb.Event) (*eventspb.EventID, error) {
	res := eventspb.EventID{}
	resp, err := doRequest(fmt.Sprintf("%s/event.EventService/CreateEvent", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func (t TestClientHTTP) UpdateEvent(_ context.Context, in *eventspb.Event) (*eventspb.Empty, error) {
	res := eventspb.Empty{}
	resp, err := doRequest(fmt.Sprintf("%s/event.EventService/UpdateEvent", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func (t TestClientHTTP) RemoveEvent(_ context.Context, in *eventspb.EventID) (*eventspb.Empty, error) {
	res := eventspb.Empty{}
	resp, err := doRequest(fmt.Sprintf("%s/event.EventService/RemoveEvent", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func (t TestClientHTTP) GetDayEvents(_ context.Context, in *eventspb.Time) (*eventspb.Events, error) {
	res := eventspb.Events{}
	resp, err := doRequest(fmt.Sprintf("%s/event.EventService/GetDayEvents", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func (t TestClientHTTP) GetWeekEvents(_ context.Context, in *eventspb.Time) (*eventspb.Events, error) {
	res := eventspb.Events{}
	resp, err := doRequest(fmt.Sprintf("%s/event.EventService/GetWeekEvents", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func (t TestClientHTTP) GetMonthEvents(_ context.Context, in *eventspb.Time) (*eventspb.Events, error) {
	res := eventspb.Events{}
	resp, err := doRequest(fmt.Sprintf("%s/event.EventService/GetMonthEvents", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func doRequest(url string, body interface{}) (*http.Response, error) {
	js, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func checkRespError(body []byte) error {
	s := spb.Status{}
	if err := json.Unmarshal(body, &s); err != nil || s.Message == "" {
		return nil
	}

	return errors.New(s.Message)
}

func getResponse(r *http.Response, v interface{}) error {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := checkRespError(body); err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func TestHTTP(t *testing.T) {
	t.Run("EmptyMemoryStorage", func(t *testing.T) { doTestMemory(t, server.EmptyStorageImpl) })
	t.Run("EmptySQLStorage", func(t *testing.T) { /*doTestSQL(t, server.EmptyStorageImpl)*/ })
	t.Run("CommonMemoryStorage", func(t *testing.T) { doTestMemory(t, server.CommonImpl) })
	t.Run("CommonSQLStorage", func(t *testing.T) { /*doTestSQL(t, server.CommonImpl)*/ })
}

func doTestMemory(t *testing.T, testFn func(t *testing.T, client server.TestClient)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := server.GetMemoryStorage()
	assert.NoError(t, err)
	a := app.NewApp(storage)

	muxMem := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/json", &runtime.JSONBuiltin{}))
	assert.NoError(t, eventspb.RegisterEventServiceHandlerServer(ctx, muxMem, a))

	srv := httptest.NewServer(muxMem)
	client := TestClientHTTP{
		addr: srv.URL,
	}

	defer srv.Close()
	testFn(t, client)
}

func doTestSQL(t *testing.T, testFn func(t *testing.T, client server.TestClient)) { //nolint:deadcode,unused
	// TODO: docker compose
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := server.GetSQLStorage()
	assert.NoError(t, err)
	a := app.NewApp(storage)
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/json", &runtime.JSONBuiltin{}))
	assert.NoError(t, eventspb.RegisterEventServiceHandlerServer(ctx, mux, a))

	srv := httptest.NewServer(mux)
	client := TestClientHTTP{
		addr: srv.URL,
	}

	defer srv.Close()
	testFn(t, client)
}
