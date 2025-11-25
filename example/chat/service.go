package chat

// ChatService is a small example interface for the codegen demo
type ChatService interface {
	Login(user string) (bool, error)
}
