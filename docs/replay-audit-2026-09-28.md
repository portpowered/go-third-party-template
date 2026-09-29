# Paired replay audit — 2026-09-28

This audit applies library standard 15 to the four standalone provider modules.
It distinguishes a full wire replay from a response body used only to test a
parser or model. Qualifying replay fixtures below are synthetic. Historical
references with unproven capture provenance are kept separate and excluded
from qualifying counts; no live account behavior is inferred from these tests.

| Library | HTTP request/response pairs | Other network exchanges | Verdict |
| --- | --- | --- | --- |
| go-ring | The implementation at `c0d4367` supplies qualifying synthetic request/response pairs for 39/39 OpenAPI operations, enforces ordered request matching, and checks socket Host/Origin. Thirty older HTTP files without capture metadata were moved to `historical/` and excluded from qualifying replay counts. | Four stored production-path synthetic transcripts map all 17 AsyncAPI channels across push, heartbeat, playback, and live/PTZ; two older signaling files without provenance are historical. | Independent LIB-05/item-15 re-audit at `c0d4367` passed all HTTP and signaling pairs, mismatch/order/consumption tests, provenance separation, local checks, and full CI. The report signoff is at `4f8cb7a`. |
| go-tuya | Synthetic pairs now cover all 28 OpenAPI HTTP operations. The ordered matcher checks the outbound request before returning the response, including encrypted query/body content and signatures. | An MQTT transcript exercises the real `MessageQueue.Start` path, including its paired HTTP configuration call and injected MQTT client. | Item 15 and the full 15-item checklist independently signed off in `docs/independent-verification.md` at v0.3.3 (`2b598fe`); exact-tag release, CI, docs, public proxy, and versioned Go Reference passed. The signoff commit `304356a` changes documentation only. |
| go-tplink | Fifteen synthetic normal fixtures pair method, origin, path, query, headers, complete request-body variants, and response. At `2298ea5`, the transport uses explicit ordered one-time steps, exact synthetic token values, and stored pairs for 33 failure/transport outcomes; it rejects mismatches and extra calls. | No additional network protocol in the current public client inventory. | Independent item-15 re-audit at `2298ea5` passed all 22 normal variants, failure pairs, mismatch/order/consumption tests, and the operation inventory. The report and checklist signoff are at `dc758c7`; CI and docs passed on the implementation commit. |
| go-alexa | Checked-in synthetic pairs cover all 18 dispatched REST routes and seven GraphQL operations. The matcher validates outbound method, full URL, headers, and body, rejects mismatches and duplicates, and asserts consumption. Response-only model inputs are kept separate. | Two event and keepalive pairs cover the supported stream boundary. | Item 15 and the full 15-item checklist independently signed off in `docs/independent-verification-final.md` at v0.3.1 (`b04ff67`); exact-tag release, CI, docs, public proxy, and versioned Go Reference passed. The signoff commit `41cc4eb` changes documentation only. |

Each of the four libraries now has an independent item-15/LIB-05 signoff for
its implemented wire surface. Alexa and Tuya also have full 15-item release
reviews at their exact tags. The pairs are synthetic compatibility evidence;
they do not replace dated provider captures or document the maintainer's
reported account tests. A passing unit test or coverage percentage alone
would not have closed this requirement.
