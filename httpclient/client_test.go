package httpclient_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	service "github.com/example/your-service-go"
	"github.com/example/your-service-go/httpclient"
)

const testBaseURL = "https://api.example.test"

type doerFunc func(*http.Request) (*http.Response, error)

func (doer doerFunc) Do(request *http.Request) (*http.Response, error) {
	return doer(request)
}

func testResponse(statusCode int, body io.ReadCloser, header http.Header) *http.Response {
	response := new(http.Response)
	response.StatusCode = statusCode
	response.Body = body
	response.Header = header

	return response
}

func widgetRequest(id string) service.GetWidgetRequest {
	return service.GetWidgetRequest{
		Auth: service.AuthContext{AccessToken: ""},
		ID:   id,
	}
}

func mustNewTestClient(t *testing.T, baseURL string, doer doerFunc) *httpclient.Client {
	t.Helper()

	client, err := httpclient.New(httpclient.Options{BaseURL: baseURL, HTTPClient: doer})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	return client
}

type contextKey struct{}

func TestGetWidgetBuildsRequestAndMapsResponse(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), contextKey{}, "request-value")
	doer := doerFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet {
			t.Errorf("method = %q, want %q", request.Method, http.MethodGet)
		}

		if request.URL.EscapedPath() != "/v1/widgets/widget%2Fone" {
			t.Errorf("path = %q, want %q", request.URL.EscapedPath(), "/v1/widgets/widget%2Fone")
		}

		if got := request.Header.Get("Authorization"); got != "Bearer example-token" {
			t.Errorf("authorization = %q, want bearer token", got)
		}

		if got := request.Context().Value(contextKey{}); got != "request-value" {
			t.Errorf("context value = %v, want request-value", got)
		}

		return testResponse(
			http.StatusOK,
			io.NopCloser(strings.NewReader("{\"id\":\"widget/one\",\"name\":\"Example\"}")),
			make(http.Header),
		), nil
	})

	client := mustNewTestClient(t, "https://api.example.test/v1", doer)

	got, err := client.GetWidget(ctx, service.GetWidgetRequest{
		Auth: service.AuthContext{AccessToken: "example-token"},
		ID:   "widget/one",
	})
	if err != nil {
		t.Fatalf("GetWidget() error = %v", err)
	}

	if got != (service.Widget{ID: "widget/one", Name: "Example"}) {
		t.Errorf("GetWidget() = %+v, want widget result", got)
	}
}

func TestGetWidgetReturnsTypedNotFoundError(t *testing.T) {
	t.Parallel()

	doer := doerFunc(func(*http.Request) (*http.Response, error) {
		return testResponse(
			http.StatusNotFound,
			io.NopCloser(strings.NewReader("")),
			http.Header{"X-Request-Id": []string{"request-123"}},
		), nil
	})
	client := mustNewTestClient(t, testBaseURL, doer)

	_, err := client.GetWidget(context.Background(), widgetRequest("missing"))

	var clientErr *service.Error
	if !errors.As(err, &clientErr) {
		t.Fatalf("GetWidget() error = %T, want *service.Error", err)
	}

	if clientErr.Kind != service.ErrorKindNotFound {
		t.Errorf("error kind = %q, want %q", clientErr.Kind, service.ErrorKindNotFound)
	}

	if clientErr.StatusCode != http.StatusNotFound {
		t.Errorf("status code = %d, want %d", clientErr.StatusCode, http.StatusNotFound)
	}

	if clientErr.RequestID != "request-123" {
		t.Errorf("request ID = %q, want request-123", clientErr.RequestID)
	}
}

func TestGetWidgetRejectsResponseWithoutBody(t *testing.T) {
	t.Parallel()

	doer := doerFunc(func(*http.Request) (*http.Response, error) {
		return testResponse(http.StatusOK, nil, nil), nil
	})
	client := mustNewTestClient(t, testBaseURL, doer)

	_, err := client.GetWidget(context.Background(), widgetRequest("widget-1"))

	var clientErr *service.Error
	if !errors.As(err, &clientErr) {
		t.Fatalf("GetWidget() error = %T, want *service.Error", err)
	}

	if clientErr.Kind != service.ErrorKindInvalidResponse {
		t.Errorf("error kind = %q, want %q", clientErr.Kind, service.ErrorKindInvalidResponse)
	}
}

func TestGetWidgetPreservesTransportCause(t *testing.T) {
	t.Parallel()

	transportCause := io.EOF
	doer := doerFunc(func(*http.Request) (*http.Response, error) {
		return nil, transportCause
	})
	client := mustNewTestClient(t, testBaseURL, doer)

	_, err := client.GetWidget(context.Background(), widgetRequest("widget-1"))

	var clientErr *service.Error
	if !errors.As(err, &clientErr) {
		t.Fatalf("GetWidget() error = %T, want *service.Error", err)
	}

	if clientErr.Kind != service.ErrorKindTransport {
		t.Errorf("error kind = %q, want %q", clientErr.Kind, service.ErrorKindTransport)
	}

	if !errors.Is(err, transportCause) {
		t.Errorf("GetWidget() error %v does not preserve transport cause", err)
	}
}
