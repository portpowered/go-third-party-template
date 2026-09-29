package service

import "strconv"

// ErrorKind identifies the class of a client failure.
type ErrorKind string

const (
	// ErrorKindConfiguration means the client could not be configured.
	ErrorKindConfiguration ErrorKind = "configuration"
	// ErrorKindInvalidRequest means the caller supplied an invalid request.
	ErrorKindInvalidRequest ErrorKind = "invalid_request"
	// ErrorKindTransport means the request could not reach the provider.
	ErrorKindTransport ErrorKind = "transport"
	// ErrorKindUnauthorized means the provider rejected the credentials.
	ErrorKindUnauthorized ErrorKind = "unauthorized"
	// ErrorKindNotFound means the requested resource does not exist.
	ErrorKindNotFound ErrorKind = "not_found"
	// ErrorKindRateLimited means the provider asked the client to slow down.
	ErrorKindRateLimited ErrorKind = "rate_limited"
	// ErrorKindServer means the provider returned a server error.
	ErrorKindServer ErrorKind = "server"
	// ErrorKindUnexpectedStatus means the provider returned another non-success status.
	ErrorKindUnexpectedStatus ErrorKind = "unexpected_status"
	// ErrorKindInvalidResponse means the response could not be decoded.
	ErrorKindInvalidResponse ErrorKind = "invalid_response"
)

// Error describes a typed client failure while preserving its underlying cause.
type Error struct {
	// Operation names the public client method that failed.
	Operation string
	// Kind classifies the failure for caller decisions.
	Kind ErrorKind
	// StatusCode is the provider's HTTP status, when available.
	StatusCode int
	// RequestID is the provider's request identifier, when available.
	RequestID string
	// Cause is the underlying error and is also available through Unwrap.
	Cause error
}

// Error returns a concise description suitable for logs.
func (err *Error) Error() string {
	if err == nil {
		return "service: <nil>"
	}

	message := "service: " + err.Operation + ": " + string(err.Kind)
	if err.StatusCode != 0 {
		message += " (HTTP " + strconv.Itoa(err.StatusCode) + ")"
	}

	if err.Cause != nil {
		message += ": " + err.Cause.Error()
	}

	return message
}

// Unwrap exposes the underlying cause for errors.Is and errors.As.
func (err *Error) Unwrap() error {
	if err == nil {
		return nil
	}

	return err.Cause
}
