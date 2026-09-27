# Go third-party service client template

This repository is a starting point for a Go client library for one service.
Replace github.com/example/your-service-go, the service package name,
SERVICE_* environment variables, and the example API URL before publishing.
The template includes an Apache-2.0 [license](LICENSE); review it for each
derived library.

## Quick start

After replacing the module path and package names:

```sh
go get github.com/example/your-service-go
```

The example reads its API base URL, account token, and widget ID from the
environment:

```sh
SERVICE_API_BASE_URL=https://api.example.com/v1 \
SERVICE_ACCESS_TOKEN=replace-me \
SERVICE_WIDGET_ID=widget-123 \
go run ./examples/basic
```

The API base URL and operation shown here are examples. Replace them with
verified service behavior before describing them as supported operations.

## Client API

The root package contains the provider-specific Client contract and public
request and result types. Its HTTP implementation is in httpclient. A
request carries its own AuthContext, so a shared client does not keep account
credentials in mutable shared state.

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "os"
    "time"

    service "github.com/example/your-service-go"
    "github.com/example/your-service-go/httpclient"
)

func main() {
    client, err := httpclient.New(httpclient.Options{
        BaseURL: "https://api.example.com/v1",
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    widget, err := client.GetWidget(ctx, service.GetWidgetRequest{
        Auth: service.AuthContext{AccessToken: os.Getenv("SERVICE_ACCESS_TOKEN")},
        ID:   "widget-123",
    })
    if err != nil {
        var clientErr *service.Error
        if errors.As(err, &clientErr) && clientErr.Kind == service.ErrorKindNotFound {
            log.Print("widget was not found")
        }
        log.Fatal(err)
    }
    fmt.Println(widget.Name)
}
```

httpclient.Options.HTTPClient accepts any implementation of HTTPDoer.
Inject an http.Client with a custom Transport for proxies, tracing, or offline
tests. The injected implementation must be safe for concurrent use when the
service client is shared.

Pass a context with a deadline to every operation. The service client owns its
HTTP dependency; callers own request contexts and any long-running session
objects they open.

## Sessions and capabilities

This template shows a request/response client. If a service also supports
long-running connections, follow the optional session guidance in
[docs/client-design.md](docs/client-design.md). Keep a session explicit,
cancellable, observable, and safe to close more than once.

The Client interface in this template describes one provider. It demonstrates
shared client conventions such as named request types, request-scoped
credentials, context cancellation, typed errors, and transport injection. It
does not define a cross-vendor interface. Add narrow capability interfaces only
after multiple independent provider implementations establish matching behavior;
[docs/client-design.md](docs/client-design.md) gives the criteria.

## Verification and releases

Run make check before proposing a change. It runs static checks, builds packages
and examples, and runs tests with the race detector. See
[docs/verification.md](docs/verification.md) for fixture provenance and CI
details, [docs/migration-guide.md](docs/migration-guide.md) for extracting a
provider client from an application, the
[Python-to-Go migration playbook](docs/python-provider-migration-playbook.md)
for a task-by-task conversion checklist, and
[docs/releasing.md](docs/releasing.md) for module tags and release verification.
