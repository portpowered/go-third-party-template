# Extracting a provider client from an application

For a task-by-task conversion of an existing Python provider library into a
standalone Go module, use the [Python-to-Go migration playbook](python-provider-migration-playbook.md).

The future go-tplink application is a useful example of this migration. The
current TP-Link implementation and tests are in
portos-backend/integrations/go-tplink. This guide describes a future extraction;
it does not claim the template is already a TP-Link client.

## Keep the provider boundary small

Move code that speaks the provider protocol into the library: request and
response encoding, endpoint or RPC details, authentication exchange and
token use, provider error decoding, resource models, and the transport adapter.
Give the public package named request and result types that describe the
provider operations callers need.

Leave application orchestration in the backend. This includes credential and
account persistence, job scheduling, application service wiring, provider
selection, and coordination with unrelated integrations. The library should
accept credentials and caller contexts at its public boundary instead of
depending on backend storage or lifecycle packages.

## Suggested extraction steps

1. Inventory the provider calls, protocol models, authentication flow, typed
   errors, tests, and fixture sources in the existing integration.
2. Mark each behavior as captured, synthetic, or historical. Keep unknown
   provenance out of the captured fixture tree.
3. Create the provider library module and move only provider protocol code and
   its focused tests. Replace backend imports with standard-library or
   provider-library packages.
4. Define a small provider-specific Client interface and named public request
   and result types. Keep account credentials on each request when the client
   may serve multiple accounts.
5. Inject HTTP or other transports, preserve context cancellation, and map
   provider failures to typed library errors.
6. Change the backend integration into an adapter that constructs the library
   client and connects its operations to backend orchestration. Keep backend
   persistence and scheduling on that side of the adapter.
7. Build, test with the race detector, vet, review the public module from a
   separate consumer, and then remove the duplicated provider implementation
   from the backend.

## Provider APIs and shared capabilities

The extracted go-tplink client should first express TP-Link's own operations.
The provider-specific client contract can share conventions with other
libraries without asserting that those libraries have identical behavior.
Introduce a cross-vendor capability interface only after independent provider
implementations validate the same inputs, result meaning, lifecycle, and error
behavior.

Do not copy an application-facing interface wholesale just because it exists.
Separate provider operations from methods added to make one backend's mocks,
configuration, or workflow convenient. Keep the adapter responsible for
translating between the backend contract and the standalone library.
