package chat

import "testing"

// TestChatTypesCompile ensures the expected generated client types are
// present at compile time (guarding the toy generator output used in
// examples).
func TestChatTypesCompile(t *testing.T) {
	// compile-time check that generated client type exists
	var _ = (*ChatClient)(nil)
}
