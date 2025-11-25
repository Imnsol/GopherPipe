package gopherpipe

// Minimal RPC envelope used for data frames. We use gob as the prototype codec.

type RPCType byte

const (
	Unary        RPCType = 0x01
	ClientStream RPCType = 0x02
	ServerStream RPCType = 0x03
	BiDi         RPCType = 0x04
)

// Envelope is the codec-serializable RPC envelope
type Envelope struct {
	RPCType     RPCType
	ServiceName string
	MethodName  string
	CallID      uint64
	Body        []byte
}
