// Package bench contains micro-benchmarks used to compare serialization
// options in the prototype (gob vs JSON for example). These benchmarks are
// intentionally lightweight and intended to be run locally in dev or CI.
package bench

import (
	"testing"

	"encoding/json"

	"github.com/anthony/gopher-pipe/internal/codec"
)

// makeProtoUser previously created a protobuf message — in this environment
// we benchmark gob vs json (protobuf requires protoc-generated types).
func makeJSONUser() GobUser {
	return GobUser{Id: 12345, Name: "Anthony", Email: "anthony@example.com"}
}

type GobUser struct {
	Id    int64
	Name  string
	Email string
}

func makeGobUser() GobUser {
	return GobUser{Id: 12345, Name: "Anthony", Email: "anthony@example.com"}
}

// Benchmark_Gob_Marshal_Unmarshal measures encoding/gob round-trip
// performance for a small user-like struct.
func Benchmark_Gob_Marshal_Unmarshal(b *testing.B) {
	u := makeGobUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf2, err := codec.Encode(u)
		if err != nil {
			b.Fatal(err)
		}
		var dest GobUser
		if err := codec.Decode(buf2, &dest); err != nil {
			b.Fatal(err)
		}
		_ = dest
	}
}

// Benchmark_JSON_Marshal_Unmarshal measures encoding/json round-trip
// performance for the same small struct so results can be compared.
func Benchmark_JSON_Marshal_Unmarshal(b *testing.B) {
	u := makeJSONUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bts, err := json.Marshal(u)
		if err != nil {
			b.Fatal(err)
		}
		var dest GobUser
		if err := json.Unmarshal(bts, &dest); err != nil {
			b.Fatal(err)
		}
		_ = dest
	}
}

// small sanity test driver for bench user creation
// TestMakeUser provides a tiny sanity check used by maintainers to
// validate the bench helper produces sensible test data.
func TestMakeUser(t *testing.T) {
	u := makeJSONUser()
	if u.Id == 0 {
		t.Fatal("id zero")
	}
	// sanity check done
}
