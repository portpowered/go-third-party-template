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
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	client, err := httpclient.New(httpclient.Options{
		BaseURL: os.Getenv("SERVICE_API_BASE_URL"),
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	widget, err := client.GetWidget(ctx, service.GetWidgetRequest{
		Auth: service.AuthContext{AccessToken: os.Getenv("SERVICE_ACCESS_TOKEN")},
		ID:   os.Getenv("SERVICE_WIDGET_ID"),
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s: %s\n", widget.ID, widget.Name)
	return nil
}
