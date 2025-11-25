# GopherPipe — Go-native RPC framework

GopherPipe is a Go‑first RPC design that rethinks remote procedure calls for Go-to-Go services.
It uses idiomatic Go primitives (channels, interfaces, goroutines) as the developer-facing API and a lightweight, optimized transport under the hood for low-latency, low-allocation RPC between Go services.

This repository contains short design notes and a starting spec for the project. Two primary design documents in this workspace are:

- `Design a competitor to gRPC using Go. Think out o....md` — overall design, contract model, transport negotiation, and trade-offs.
- `Expand on gopherpipe.md` — deeper notes about TCP_LITE framing, discovery (GopherMap), zero-copy/gob optimizations, and usage examples.

---

## Quick summary ✅

- Philosophy: Prioritize Go-to-Go performance and developer ergonomics over polyglot/browser compatibility.
- API model: Define service contracts as Go interfaces. Streaming and RPCs expose typed channels to the user (chan-based streaming).
- Default codec: `encoding/gob` (pluggable codec architecture for Protobuf/JSON/MessagePack when interoperability is needed).
- Transport: `TCP_LITE` — a minimal, length-prefixed framing protocol geared toward zero-allocation parsing. Optional alternative transports (WebSockets / UNIX sockets) supported via negotiation.
- Discovery: Built-in gossip (GopherMap) for zero-config peer discovery.

---

## Why build GopherPipe? 💡

- Remove unnecessary complexity for internal Go microservices (no .proto files forced).
- Make streaming and RPC feel like native Go concurrency — use channels and select statements instead of generated Send/Recv loops.
- Optimize for low-latency, low-GC pressure with zero-copy and io.ReaderFrom/io.WriterTo optimizations where possible.

---

## Features & tradeoffs

- Features:
  - Channel-first streaming API and idiomatic Go contracts
  - Pluggable codec layer (gob by default) for fast Go-to-Go communication
  - Lightweight transport (`TCP_LITE`) with a 5-byte header option for minimal parsing
  - Mesh-friendly built-in gossip for discovery

- Tradeoffs:
  - Default settings favor Go-only environments — browser support and other languages are not first-class without switching codecs/transports
  - Requires additional ecosystem tooling for advanced features gRPC enjoys out of the box (xDS, envoy integration, extensive multi-language clients)

---

## Example (from the design notes)

This is a simplified client/server example demonstrating the channel-first, idiomatic usage pattern described in the docs.

Server-side (conceptual):

```go
// type ChatService interface { 
//     JoinRoom(roomID string, incoming <-chan string) (<-chan Message, error)
// }

func (s *Server) JoinRoom(roomID string, incoming <-chan string) (<-chan Message, error) {
    outgoing := make(chan Message)
    go func() {
        defer close(outgoing)
        for text := range incoming {
            msg := Message{Sender: "System", Text: "Echo: " + text}
            outgoing <- msg
        }
    }()
    return outgoing, nil
}
```

Client-side (conceptual):

```go
input := make(chan string)
output, _ := client.JoinRoom("General", input)

go func() {
    input <- "Hello Gophers!"
    input <- "Is this thing on?"
    close(input)
}()

for msg := range output {
    fmt.Printf("[%s]: %s\n", msg.Sender, msg.Text)
}
```

---

## Project status & quick tour 🔭

Progress: I implemented a working prototype and scaffolded a small proof-of-concept library plus tests and a design RFC. The repo now contains a runnable tiny `TCP_LITE` echo server + gob client prototype, the RFC describing negotiation and wire-format, a minimal `gopherpipe` package (client/server/envelope), a toy codegen tool and example Chat service + generated client stub.

What's currently included (high level):

- Prototype: `internal/tcplite` framing, `internal/codec` (gob), `internal/server` + `cmd/echoserver` and `cmd/echoclient` with unit tests.
- RFC: `RFC-TCP_LITE.md` documents negotiation, frame formats and service registration semantics.
- PoC library: `gopherpipe/` contains a minimal Envelope API, client/server prototypes and reflection-based dispatch used for examples.
- Example service & codegen: `example/chat/service.go`, `cmd/genstub` and generated `example/chat/client_gen.go` plus example CLI under `example/chatcmd/`.
- Defensive handling: the server now detects non-TCP_LITE (e.g., accidental HTTP probes) and replies with a friendly HTTP 400 (tests included).

---

## How to run (quick)

Start the TCP_LITE echo server:

```pwsh
cd 'C:\Users\Anthony\Desktop\gopher-pipe'
go run ./cmd/echoserver
```

Open another terminal and run the client that sends one message and reads the echo back:

```pwsh
cd 'C:\Users\Anthony\Desktop\gopher-pipe'
go run ./cmd/echoclient
```

Example Chat service / client (small PoC):

```pwsh
# Terminal A: run example server
go run ./example/chatcmd/server

# Terminal B: run example client
go run ./example/chatcmd/client
```

---

## Tests & Benchmarks ✅

This repository contains several kinds of tests. Below are recommended commands and where those tests live:

- Unit tests — quick, package-level tests that exercise small components and logic.
    - Where: spread across packages such as `internal/codec`, `internal/message`, `gopherpipe`, and many `*_test.go` files throughout the workspace.
    - Run all unit tests:

```pwsh
go test ./...
```

- Build tests — lightweight checks that ensure important build targets and example programs compile correctly. These typically live as `build_test.go` files under `cmd/*` and `genstub/`.
    - Where: `cmd/echoclient/build_test.go`, `cmd/echoserver/build_test.go`, `cmd/genstub/build_test.go`, etc.
    - Run only build tests (match by name):

```pwsh
# run tests with "build" in their name
go test ./... -run Build
```

- Integration tests — run end-to-end flows and multi-component behaviors (may require network or ports and are slower).
    - Where: `internal/server/integration_test.go` and other `integration_test.go` files across packages.
    - Run a specific integration test file / package:

```pwsh
# run only the integration tests in internal/server
go test ./internal/server -run TestIntegration -v
```

- Full test suite (unit + integration + build checks):

```pwsh
go test ./...
```

- Benchmarks (serialization) — benchmarks live under `/bench` and package-level benchmarks in packages when present:

```pwsh
go test ./bench -bench . -run ^$
```

Notes & tips:

- If you only want unit tests and want to skip long-running integration tests, consider running the package list you care about directly, e.g.: `go test ./internal/codec ./message`.
- Integration tests may depend on network resources or ports — run them in a controlled environment or CI with reproducible ports, and add `-v` to see detailed logs.
- Build tests that live in `build_test.go` often validate `go build` or `go list` and are designed to quickly catch compile-time regressions for example CLIs or generated stubs.

Note: a true gob vs protobuf benchmark requires protobuf message types generated by `protoc` (or committed pre-generated `.pb.go`). Protobuf toolchain (protoc + protoc-gen-go) isn't available by default in this environment — you can either install `protoc` and run the `proto/` generation steps locally, or I can add pre-generated `.pb.go` files into the repo to enable accurate protobuf benchmarks.

---

## Contributing & next steps

If you'd like me to continue, here are suggested next actions (pick one):

1. Full streaming prototype — implement channel-backed streaming (client + server) and end-to-end examples.
2. Protobuf benches — add pre-generated `.pb.go` files or run `protoc` locally then rerun benchmarks for accurate gob vs protobuf comparisons.
3. Codegen & library — expand the code generation into a polished `gopherpipe` tool that emits client/server stubs from Go interfaces and add integration tests.
4. RFC -> spec expansion — finalize the protocol RFC (detailed envelopes, TLV options, security/negotiation details and diagrams).

Tell me which one to focus on next and I'll implement it or provide a detailed plan and prototype.

---

Thank you for the design notes — this README aims to centralize them and provide a clear path for prototyping and iteration.
