// Package httpclient implements the service client over HTTP.
package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	service "github.com/example/your-service-go"
)

const maxResponseBytes = 1 << 20

const getWidgetOperation = "GetWidget"

// HTTPDoer is the part of net/http.Client used by this package.
//
// An injected implementation must support concurrent calls if the client is
// shared between goroutines.
type HTTPDoer interface {
	// Do sends an HTTP request and returns the provider response.
	Do(request *http.Request) (*http.Response, error)
}

// Options configures the HTTP client.
type Options struct {
	// BaseURL is the provider API root, such as "https://api.example.com/v1".
	BaseURL string
	// HTTPClient replaces the default http.DefaultClient.
	HTTPClient HTTPDoer
}

// Client implements service.Client using HTTP.
type Client struct {
	baseURL    *url.URL
	httpClient HTTPDoer
}

// New creates an HTTP client. The caller should use request contexts with
// deadlines for network operations.
func New(options Options) (*Client, error) {
	baseURL, err := url.Parse(options.BaseURL)
	if err != nil {
		return nil, clientError("New", service.ErrorKindConfiguration, 0, "", err)
	}

	if (baseURL.Scheme != "https" && baseURL.Scheme != "http") ||
		baseURL.Host == "" || baseURL.User != nil || baseURL.RawQuery != "" ||
		baseURL.Fragment != "" {
		return nil, clientError("New", service.ErrorKindConfiguration, 0, "", nil)
	}

	httpClient := options.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{baseURL: baseURL, httpClient: httpClient}, nil
}

// GetWidget fetches one widget using the request's credentials.
func (client *Client) GetWidget(
	ctx context.Context,
	request service.GetWidgetRequest,
) (service.Widget, error) {
	if ctx == nil || request.ID == "" {
		return service.Widget{}, clientError(getWidgetOperation, service.ErrorKindInvalidRequest, 0, "", nil)
	}

	endpoint, err := widgetURL(client.baseURL, request.ID)
	if err != nil {
		return service.Widget{}, clientError(getWidgetOperation, service.ErrorKindInvalidRequest, 0, "", err)
	}

	httpRequest, err := newWidgetRequest(ctx, endpoint, request.Auth.AccessToken)
	if err != nil {
		return service.Widget{}, clientError(getWidgetOperation, service.ErrorKindInvalidRequest, 0, "", err)
	}

	response, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return service.Widget{}, clientError(getWidgetOperation, service.ErrorKindTransport, 0, "", err)
	}

	if response == nil {
		return service.Widget{}, clientError(getWidgetOperation, service.ErrorKindTransport, 0, "", nil)
	}

	if response.Body != nil {
		defer func() { _ = response.Body.Close() }()
	}

	return decodeWidgetResponse(response)
}

func newWidgetRequest(ctx context.Context, endpoint *url.URL, accessToken string) (*http.Request, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create widget request: %w", err)
	}

	httpRequest.Header.Set("Accept", "application/json")

	if accessToken != "" {
		httpRequest.Header.Set("Authorization", "Bearer "+accessToken)
	}

	return httpRequest, nil
}

func decodeWidgetResponse(response *http.Response) (service.Widget, error) {
	if response.Body == nil {
		return service.Widget{}, clientError(getWidgetOperation, service.ErrorKindInvalidResponse, 0, "", nil)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return service.Widget{}, statusError(getWidgetOperation, response)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil {
		return service.Widget{}, clientError(getWidgetOperation, service.ErrorKindInvalidResponse, 0, "", err)
	}

	if len(body) > maxResponseBytes {
		return service.Widget{}, clientError(getWidgetOperation, service.ErrorKindInvalidResponse, 0, "", nil)
	}

	return decodeWidget(body)
}

func decodeWidget(body []byte) (service.Widget, error) {
	var wire widgetResponse

	err := json.Unmarshal(body, &wire)
	if err != nil {
		return service.Widget{}, clientError(getWidgetOperation, service.ErrorKindInvalidResponse, 0, "", err)
	}

	if wire.ID == "" {
		return service.Widget{}, clientError(getWidgetOperation, service.ErrorKindInvalidResponse, 0, "", nil)
	}

	return service.Widget{ID: wire.ID, Name: wire.Name}, nil
}

type widgetResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func widgetURL(baseURL *url.URL, id string) (*url.URL, error) {
	endpoint := *baseURL
	rawPath := strings.TrimRight(baseURL.EscapedPath(), "/") + "/widgets/" + url.PathEscape(id)

	path, err := url.PathUnescape(rawPath)
	if err != nil {
		return nil, fmt.Errorf("unescape widget URL path: %w", err)
	}

	endpoint.Path = path
	endpoint.RawPath = rawPath

	return &endpoint, nil
}

func statusError(operation string, response *http.Response) *service.Error {
	kind := service.ErrorKindUnexpectedStatus

	switch response.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		kind = service.ErrorKindUnauthorized
	case http.StatusNotFound:
		kind = service.ErrorKindNotFound
	case http.StatusTooManyRequests:
		kind = service.ErrorKindRateLimited
	default:
		if response.StatusCode >= http.StatusInternalServerError {
			kind = service.ErrorKindServer
		}
	}

	return clientError(operation, kind, response.StatusCode, response.Header.Get("X-Request-ID"), nil)
}

func clientError(
	operation string,
	kind service.ErrorKind,
	statusCode int,
	requestID string,
	cause error,
) *service.Error {
	return &service.Error{
		Operation:  operation,
		Kind:       kind,
		StatusCode: statusCode,
		RequestID:  requestID,
		Cause:      cause,
	}
}
