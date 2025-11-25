// Package tcplite implements a compact, length-prefixed wire framing
// protocol used by the gopherpipe prototype. The framing is intentionally
// small: a single-byte type followed by a 4-byte big-endian length and the
// payload. The implementation focuses on predictable parsing and a low
// allocation path for small messages.
package tcplite

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// TCP_LITE frame header: 1 byte Type, 4 bytes Length (big endian)
// Frame type constants used on the wire for TCP_LITE frames.
const (
	FrameTypeData      byte = 0x01
	FrameTypeHeartbeat byte = 0x02
	FrameTypeError     byte = 0x03
	FrameTypeClose     byte = 0x04
	// optional future frame types
	FrameTypeServiceReg    byte = 0x05
	FrameTypeServiceLookup byte = 0x06
)

// WriteFrame writes a TCP_LITE frame to w using the canonical header
// layout (1-byte type, 4-byte length big-endian, payload). It returns any
// write error from the underlying writer.
func WriteFrame(w io.Writer, ftype byte, payload []byte) error {
	header := make([]byte, 5)
	header[0] = ftype
	binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))
	if _, err := w.Write(header); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

// ReadFrame reads a single frame from r and returns the frame type and
// payload. The function performs basic validation of frame types and caps
// payload length with a sanity check to prevent large/allocation attacks.
func ReadFrame(r io.Reader) (byte, []byte, error) {
	header := make([]byte, 5)
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, nil, err
	}
	ftype := header[0]

	// validate frame type before trusting length bytes — if someone connects with HTTP
	// or another protocol we should reject early (the length bytes would otherwise look huge).
	switch ftype {
	case FrameTypeData, FrameTypeHeartbeat, FrameTypeError, FrameTypeClose, FrameTypeServiceReg, FrameTypeServiceLookup:
		// valid
	default:
		return 0, nil, &InvalidFrameHeaderError{Header: append([]byte(nil), header...)}
	}
	length := binary.BigEndian.Uint32(header[1:])
	if length > 10<<20 { // 10MB sanity check
		return 0, nil, fmt.Errorf("frame too large: %d (header=%x)", length, header)
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return 0, nil, err
	}
	return ftype, payload, nil
}

// InvalidFrameHeaderError indicates the first header bytes did not match a known TCP_LITE frame type
// InvalidFrameHeaderError is returned when the first header byte does not
// match any known TCP_LITE frame type. The Header field contains the raw
// header bytes that failed validation.
type InvalidFrameHeaderError struct {
	Header []byte
}

// Error formats a short human-readable description for invalid header
// errors including the raw header bytes (hex) when available.
func (e *InvalidFrameHeaderError) Error() string {
	if len(e.Header) >= 5 {
		return fmt.Sprintf("invalid frame type: header=%x", e.Header[:5])
	}
	return fmt.Sprintf("invalid frame header: header=%x", e.Header)
}

// IsInvalidFrameHeader tells whether err is an InvalidFrameHeaderError
// IsInvalidFrameHeader reports whether err is an InvalidFrameHeaderError
// so callers can take specific recovery actions (for example, respond
// with an HTTP 400 to accidental HTTP probes).
func IsInvalidFrameHeader(err error) bool {
	var ie *InvalidFrameHeaderError
	return errors.As(err, &ie)
}
