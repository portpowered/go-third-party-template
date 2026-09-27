# Client design guidance

## Provider client conventions

The root service package is the public API for one provider. Keep its request
and result types independent of wire-format structs. Group the supported
operations in the Client interface so callers can see the library surface in
one place.

Use context.Context on every operation that can wait on I/O. Keep
account-specific credentials on each account request so one client can serve
multiple accounts without changing shared authorization state. Return typed
errors that preserve causes and expose a stable failure kind.

The example Client interface is a provider-specific contract. These design
conventions can be reused by separate provider libraries, but this template has
one implementation and does not validate a shared interface across providers.

## Optional cross-vendor capabilities

Add a cross-vendor capability interface only after at least two independent
provider implementations demonstrate the same operation semantics and failure
behavior. Keep each interface narrow and based on what a client can do, such as
recording retrieval or device control. Do not make every provider satisfy an
umbrella interface when support varies.

Before adding a shared interface, compare the request inputs, result meaning,
authentication needs, cancellation behavior, and error classes across the
implementations. If those differ, keep the APIs provider-specific or define a
smaller capability with behavior both implementations can honor. This template
contains no cross-vendor capability interface.

## Optional long-running sessions

Skip session types when operations are short request/response calls. When a
service operation opens a long-running connection, return an explicit session
object from the client. Pass a context when opening it, document who owns
Close, make repeated close calls safe, and surface terminal errors and
backpressure to the caller. Keep connection state in the session rather than
the reusable client.

## Transport ownership

The httpclient package owns request construction, response decoding, and
transport error classification. It accepts an injected HTTPDoer while keeping
wire structs private. The application owns its HTTP transport configuration
and request deadlines. If a client later adds another protocol, put that
transport in a focused package and keep transport details out of the root API
types.
