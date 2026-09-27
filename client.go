// Package service defines the provider-specific client contract and public
// request and result types.
package service

import "context"

// Client contains the supported operations for one provider.
//
// Implementations should be safe for concurrent use. Account credentials
// belong on each request so one client can serve multiple accounts.
type Client interface {
	GetWidget(ctx context.Context, request GetWidgetRequest) (Widget, error)
}

// AuthContext holds credentials for one account request.
type AuthContext struct {
	// AccessToken is the bearer token for this account request.
	AccessToken string
}

// GetWidgetRequest identifies a widget and the account used to access it.
type GetWidgetRequest struct {
	// Auth supplies credentials for this request.
	Auth AuthContext
	// ID is the provider's widget identifier.
	ID string
}

// Widget is the public result returned by GetWidget.
type Widget struct {
	// ID is the provider's widget identifier.
	ID string
	// Name is the display name returned by the provider.
	Name string
}
