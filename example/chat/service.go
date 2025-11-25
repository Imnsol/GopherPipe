package chat

// Package chat contains a tiny example interface used by the codegen
// demonstration. The interface is intentionally minimal for test and demo
// purposes.
// ChatService is a small example interface for the codegen demo.
type ChatService interface {
	Login(user string) (bool, error)
}
