# Paired replay audit — 2026-09-28

This audit applies library standard 15 to the four standalone provider modules.
It distinguishes a full wire replay from a response body used only to test a
parser or model. All fixtures discussed below are synthetic unless explicitly
labelled captured; no live account behavior is inferred from them.

| Library | HTTP request/response pairs | Other network exchanges | Verdict |
| --- | --- | --- | --- |
| go-ring | Captured and reference HTTP replay cassettes store method, origin, path, repeated query, headers, body, and response; `internal/testkit/replay` matches requests and asserts consumption. Older `http/baseline` bodies are historical model inputs, not HTTP replay cases. | Captured signaling fixtures store ordered client and server frames; the script asserts completion. | Paired replay structure present for the recorded routes. Audit the coverage of every supported operation before a complete item-15 sign-off. |
| go-tuya | Seven synthetic HTTP replay fixtures now store request and response together. The transport checks the request before returning the response and rejects a repeated route; tests assert call order. Encrypted query fields have format checks and decoded-value assertions. | MQTT is exercised through injected test seams, but a complete ordered publish/subscribe and event transcript has not been established in this audit. | Open: expand paired replay to the remaining HTTP operations and MQTT lifecycle. |
| go-tplink | Fifteen synthetic response fixtures now contain method, origin, path, query policy, relevant headers, permitted operation, and response. The transport rejects a request outside this envelope. Several tests separately inspect encoded command fields. | No additional network protocol in the current public client inventory. | Open: move complete command bodies and expected call count/order into replay pairs; remove permissive default routing. |
| go-alexa | The generic capture replay helper now matches method, full URL, headers, and body in order, rejects mismatches and duplicate calls, and exposes `AssertConsumed`. Its synthetic `pkg/alexa/testdata` JSON files remain response-only model inputs, while many REST and GraphQL tests use ad hoc response transports. | Event stream tests use synthetic responses and request assertions, without an exhaustive ordered bidirectional transcript. | Open: build operation-level paired replay for REST, GraphQL, and events, then re-run independent review. |

For each library, a completed verdict requires an endpoint and event inventory
mapped to paired replay cases, negative mismatch tests, consumed-call checks,
and an independent review of the final commit. A passing unit test or coverage
percentage alone does not close this requirement.
