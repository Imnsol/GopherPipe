package codec

import (
	"reflect"
	"testing"
)

type testStruct struct {
	A int
	B string
}

func TestEncodeDecodeStruct(t *testing.T) {
	in := testStruct{A: 42, B: "hello"}
	b, err := Encode(in)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var out testStruct
	if err := Decode(b, &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("mismatch: got %+v want %+v", out, in)
	}
}
