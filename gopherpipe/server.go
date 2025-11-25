package gopherpipe

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"

	"github.com/anthony/gopher-pipe/internal/codec"
	"github.com/anthony/gopher-pipe/internal/tcplite"
)

// Server is a tiny registry + TCP server used by examples. It supports
// basic service registration and a simple reflection-based unary call
// dispatcher used only for the prototype.
type Server struct {
	addr     string
	mu       sync.RWMutex
	services map[string]interface{}
}

// NewServer creates a new Server listening on the supplied address.
// The server automatically registers the Envelope type with gob so tests
// and examples can rely on stable serialization.
func NewServer(addr string) *Server {
	// register envelope type
	gob.Register(Envelope{})
	return &Server{addr: addr, services: make(map[string]interface{})}
}

// Register adds a service implementation under a logical name. Example
// code calls methods on the registered implementation via reflection.
func (s *Server) Register(name string, impl interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.services[name] = impl
}

// Serve begins listening for TCP connections and handles incoming frames
// concurrently. Serve blocks while the listener is active.
func (s *Server) Serve() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go s.handleConn(conn)
	}
}

// handleConn reads frames from a single connection and dispatches requests
// to registered services. It's invoked in a goroutine per accepted
// connection.
func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		ftype, payload, err := tcplite.ReadFrame(conn)
		if err != nil {
			log.Println("read frame error:", err)
			return
		}
		if ftype != tcplite.FrameTypeData {
			continue
		}
		var env Envelope
		if err := codec.Decode(payload, &env); err != nil {
			log.Println("decode envelope:", err)
			_ = tcplite.WriteFrame(conn, tcplite.FrameTypeError, []byte(err.Error()))
			continue
		}
		// naive: look up service and reflect-call method name if possible
		s.mu.RLock()
		impl := s.services[env.ServiceName]
		s.mu.RUnlock()
		if impl == nil {
			_ = tcplite.WriteFrame(conn, tcplite.FrameTypeError, []byte("service not found"))
			continue
		}
		// For the prototype we expect a unary call where method takes (in) and returns (out, error)
		// We'll use reflection to call the method
		respEnv, err := s.handleUnaryCall(impl, env)
		if err != nil {
			_ = tcplite.WriteFrame(conn, tcplite.FrameTypeError, []byte(err.Error()))
			continue
		}
		if err := tcplite.WriteFrame(conn, tcplite.FrameTypeData, respEnv); err != nil {
			log.Println("write reply error:", err)
			return
		}
	}
}

// handleUnaryCall performs reflection-based invocation of a simple unary
// method (one argument, returns (T, error)). It decodes the incoming
// argument, calls the method, and re-encodes the return value into a new
// Envelope payload.
func (s *Server) handleUnaryCall(impl interface{}, env Envelope) ([]byte, error) {
	// Reflection-based invocation for simple signatures.
	mv := reflect.ValueOf(impl)
	method := mv.MethodByName(env.MethodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %s not found", env.MethodName)
	}
	mtype := method.Type()
	if mtype.NumIn() != 1 {
		return nil, errors.New("only single-arg unary methods supported in prototype")
	}
	if mtype.NumOut() != 2 {
		return nil, errors.New("method must return (T, error)")
	}

	// prepare argument value of required type
	argType := mtype.In(0)
	argPtr := reflect.New(argType)
	// decode body into argPtr.Interface()
	if err := codec.Decode(env.Body, argPtr.Interface()); err != nil {
		return nil, err
	}
	args := []reflect.Value{argPtr.Elem()}
	// call method
	results := method.Call(args)
	// result value and error
	resVal := results[0].Interface()
	var callErr error
	if !results[1].IsNil() {
		callErr = results[1].Interface().(error)
	}
	if callErr != nil {
		return nil, callErr
	}
	// encode response body
	outb, err := codec.Encode(resVal)
	if err != nil {
		return nil, err
	}
	respEnv := Envelope{RPCType: Unary, ServiceName: env.ServiceName, MethodName: env.MethodName, CallID: env.CallID, Body: outb}
	return codec.Encode(respEnv)
}

func init() {
	// register common types for gob across the prototype
	codec.Encode(struct{}{})
}
