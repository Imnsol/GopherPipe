package gopherpipe

// Package-level notes: Envelope is the minimal serializable RPC envelope used
// when transporting messages between a gopherpipe client and server. The
// prototype uses encoding/gob as a convenient serialization format.

// RPCType describes the logical RPC stream type carried by an Envelope.
type RPCType byte

const (
	Unary        RPCType = 0x01
	ClientStream RPCType = 0x02
	ServerStream RPCType = 0x03
	BiDi         RPCType = 0x04
)

// Envelope is the serializable RPC envelope used on the wire. It contains
// a small set of fields sufficient for prototype unary and streaming
// operations: the RPC type, service and method names for dispatch, a
// unique CallID for matching requests/responses, and the message body.
type Envelope struct {
	RPCType     RPCType
	ServiceName string
	MethodName  string
	CallID      uint64
	Body        []byte
}
