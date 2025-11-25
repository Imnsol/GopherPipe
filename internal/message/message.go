package message

// Package message defines the small Message type used in the echo server
// examples. This keeps example payloads small and easy to encode for the
// prototype tests and demos.
// Message is a simple struct for the prototype.
type Message struct {
	ID   int64
	From string
	Body string
}
