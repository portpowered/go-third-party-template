// Package main shows a basic service client call.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	service "github.com/example/your-service-go"
	"github.com/example/your-service-go/httpclient"
)

const requestTimeout = 15 * time.Second

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	client, err := httpclient.New(httpclient.Options{
		BaseURL:    os.Getenv("SERVICE_API_BASE_URL"),
		HTTPClient: nil,
	})
	if err != nil {
		return fmt.Errorf("create service client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	widget, err := client.GetWidget(ctx, service.GetWidgetRequest{
		Auth: service.AuthContext{AccessToken: os.Getenv("SERVICE_ACCESS_TOKEN")},
		ID:   os.Getenv("SERVICE_WIDGET_ID"),
	})
	if err != nil {
		return fmt.Errorf("get widget: %w", err)
	}

	_, err = fmt.Fprintf(os.Stdout, "%s: %s\n", widget.ID, widget.Name)
	if err != nil {
		return fmt.Errorf("write widget output: %w", err)
	}

	return nil
}
