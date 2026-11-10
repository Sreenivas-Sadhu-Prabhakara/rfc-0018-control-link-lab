# RFC 18 — Control link lab

**Covered RFC 18. Next RFC 19.** Before reuse, the RFC 17 README, command,
module, license, and CI workflow were read from its public main branch at
[`d929377`](https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0017-link-expiry-lab/commit/d9293773059ed0862195ef17da51b3c10793c37c).
Number 14 is recorded as Not Issued in the RFC Editor index.

## The proposal

In September 1969, Vint Cerf at UCLA proposed reserving logical link 1 for
host-to-host control messages and link 0 for control messages between
ARPANET Interface Message Processors (IMPs). Separating these conversations
was meant to reduce interference and delays for host control messages.
This is a brief historical design proposal.

The original category field is blank. The RFC Editor catalogs it in the
Legacy stream with publication and current status Unknown. Its updates,
updated-by, obsoletes, and obsoleted-by lists are empty. Verified
10 November 2026 against the official text and metadata.

- [RFC Editor text](https://www.rfc-editor.org/rfc/rfc18.txt)
- [RFC Editor metadata](https://www.rfc-editor.org/rfc/rfc18.json)
- [RFC Editor index](https://www.rfc-editor.org/rfc-index.txt)

## Applications

1. **Network control traffic, selected.** Compare the queue delay a short
   host-control burst sees during a larger switch-control burst.
2. **Banking integration.** Model whether a reserved consumer slot could
   protect time-sensitive reconciliation commands during batch replays.
3. **Cloud operations.** Compare shared and reserved worker capacity when
   health probes and configuration reloads arrive together.

## Run the lab

Use Go 1.22 or newer from the repository root.

```sh
go run ./cmd/controllab < examples/burst.json
go run ./cmd/controllab -html < examples/burst.json > docs/index.html
go test -race ./...
go vet ./...
go build ./...
```

The JSON and HTML report the per-message queue wait and per-class maximum.
The included burst yields a host maximum of **3 → 1 ticks** and an IMP
maximum of **2 → 5 ticks** under shared FIFO and reserved slots, respectively.

## Model and tradeoff

Both policies have two unit-sized service slots per tick. Shared FIFO
can spend both on either class. Reserved slots give each class one slot
per tick, with no borrowing. Each message consumes one slot in the tick
it starts, so queue wait is start tick minus arrival tick. A zero means
immediate service. Same-tick arrivals enter each queue in input order.

Reserved slots prevent the other class from consuming your capacity,
while shared slots let a single busy class use both. The workload is a
synthetic example with logical time, equal service cost, and queues large
enough to hold the input. The model does not recreate the ARPANET wire
protocol or predict production latency. GitHub Pages displays the saved
Go-generated trace. The simulator runs locally.

The CLI accepts one JSON object with 1–500 arrivals, nondecreasing ticks
from 0 to 10000, and class `host` or `imp`. It rejects unknown fields,
trailing JSON, and input larger than 1 MiB. With an absent class, its
summary count is zero and the HTML displays a dash for maximum wait.

Tests cover exact waits, tie order, the cost of idle reserved capacity,
time gaps, malformed input, and JSON and HTML output. CI regenerates the
published trace and checks it against the committed page. MIT licensed.
