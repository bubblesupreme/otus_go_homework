package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"testing"

	eventspb "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/api"

	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/server"
	"github.com/stretchr/testify/assert"
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

const (
	host = "127.0.0.1"
	port = 8080
)

type TestClientHTTP struct {
	addr string
}

func (t TestClientHTTP) CreateEvent(_ context.Context, in *eventspb.Event) (*eventspb.EventID, error) {
	res := eventspb.EventID{}
	resp, err := doRequest(fmt.Sprintf("http://%s/event.EventService/CreateEvent", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func (t TestClientHTTP) UpdateEvent(_ context.Context, in *eventspb.Event) (*eventspb.Empty, error) {
	res := eventspb.Empty{}
	resp, err := doRequest(fmt.Sprintf("http://%s/event.EventService/UpdateEvent", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func (t TestClientHTTP) RemoveEvent(_ context.Context, in *eventspb.EventID) (*eventspb.Empty, error) {
	res := eventspb.Empty{}
	resp, err := doRequest(fmt.Sprintf("http://%s/event.EventService/RemoveEvent", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func (t TestClientHTTP) GetDayEvents(_ context.Context, in *eventspb.Time) (*eventspb.Events, error) {
	res := eventspb.Events{}
	resp, err := doRequest(fmt.Sprintf("http://%s/event.EventService/GetDayEvents", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func (t TestClientHTTP) GetWeekEvents(_ context.Context, in *eventspb.Time) (*eventspb.Events, error) {
	res := eventspb.Events{}
	resp, err := doRequest(fmt.Sprintf("http://%s/event.EventService/GetWeekEvents", t.addr), in)
	if err != nil {
		return &res, err
	}

	return &res, getResponse(resp, &res)
}

func (t TestClientHTTP) GetMonthEvents(_ context.Context, in *eventspb.Time) (*eventspb.Events, error) {
	res := eventspb.Events{}
	resp, err := doRequest(fmt.Sprintf("http://%s/event.EventService/GetMonthEvents", t.addr), in)
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

func doTest(t *testing.T, testFn func(t *testing.T, server server.Server, client server.TestClient)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := TestClientHTTP{
		addr: net.JoinHostPort(host, strconv.Itoa(port)),
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
