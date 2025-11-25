// Package codec provides a small wrapper around encoding/gob used by the
// prototypes in this repository. It exposes simple Encode/Decode helpers
// that produce/consume byte slices for easy transport over tcplite frames.
package codec

import (
	"bytes"
	"encoding/gob"
)

// Encode serializes v using encoding/gob and returns bytes
// Encode serializes the provided value using encoding/gob into a byte
// slice suitable for sending over the wire.
func Encode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode deserializes data into v (pointer)
// Decode deserializes the byte slice into v, which must be a pointer.
func Decode(data []byte, v interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(v)
}
