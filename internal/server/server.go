// Package server contains a tiny example echo server/handler used by the
// command line examples (echoserver/echoclient). It shows how to use the
// tcplite framing and the codec to decode/encode simple messages.
package server

import (
	"log"
	"net"

	"github.com/anthony/gopher-pipe/internal/codec"
	"github.com/anthony/gopher-pipe/internal/message"
	"github.com/anthony/gopher-pipe/internal/tcplite"
)

// Serve listens on addr and serves a very small echo protocol over
// tcplite frames. Each incoming data frame is expected to contain a
// message.Message marshaled with the codec package; the server echoes the
// same payload back in a data frame.
func Serve(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConn(conn)
	}
}

// handleConn drives the lifecycle for a single connection — it reads
// frames, decodes/encodes messages and handles simple control frame types
// such as close and heartbeat.
func handleConn(conn net.Conn) {
	defer conn.Close()
	log.Println("client connected:", conn.RemoteAddr())
	for {
		ftype, payload, err := tcplite.ReadFrame(conn)
		if err != nil {
			// special-case: if the header looked like HTTP (invalid frame header), reply with a friendly HTTP 400
			if tcplite.IsInvalidFrameHeader(err) {
				// read-only check for common HTTP methods in the raw header bytes
				// we don't have raw header here; the InvalidFrameHeaderError includes it in the error string
				// Try to write a minimal HTTP 400 response to help misguided clients / probes
				httpResp := "HTTP/1.1 400 Bad Request\r\nContent-Length: 11\r\nContent-Type: text/plain\r\n\r\nBad Request"
				conn.Write([]byte(httpResp))
				return
			}
			log.Println("read frame error:", err)
			return
		}
		switch ftype {
		case tcplite.FrameTypeData:
			var msg message.Message
			if err := codec.Decode(payload, &msg); err != nil {
				log.Println("decode error:", err)
				_ = tcplite.WriteFrame(conn, tcplite.FrameTypeError, []byte(err.Error()))
				continue
			}
			log.Printf("received: %+v\n", msg)
			// echo back
			resp, err := codec.Encode(msg)
			if err != nil {
				log.Println("encode error:", err)
				_ = tcplite.WriteFrame(conn, tcplite.FrameTypeError, []byte(err.Error()))
				continue
			}
			if err := tcplite.WriteFrame(conn, tcplite.FrameTypeData, resp); err != nil {
				log.Println("write frame error:", err)
				return
			}
		case tcplite.FrameTypeClose:
			log.Println("client requested close")
			return
		case tcplite.FrameTypeHeartbeat:
			// ignore
		default:
			log.Println("unknown frame type:", ftype)
		}
	}
}
