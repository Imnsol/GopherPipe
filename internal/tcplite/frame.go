package tcplite

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// TCP_LITE frame header: 1 byte Type, 4 bytes Length (big endian)
const (
	FrameTypeData      byte = 0x01
	FrameTypeHeartbeat byte = 0x02
	FrameTypeError     byte = 0x03
	FrameTypeClose     byte = 0x04
	// optional future frame types
	FrameTypeServiceReg    byte = 0x05
	FrameTypeServiceLookup byte = 0x06
)

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
type InvalidFrameHeaderError struct {
	Header []byte
}

func (e *InvalidFrameHeaderError) Error() string {
	if len(e.Header) >= 5 {
		return fmt.Sprintf("invalid frame type: header=%x", e.Header[:5])
	}
	return fmt.Sprintf("invalid frame header: header=%x", e.Header)
}

// IsInvalidFrameHeader tells whether err is an InvalidFrameHeaderError
func IsInvalidFrameHeader(err error) bool {
	var ie *InvalidFrameHeaderError
	return errors.As(err, &ie)
}
