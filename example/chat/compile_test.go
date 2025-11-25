package chat

import "testing"

func TestChatTypesCompile(t *testing.T) {
	// compile-time check that generated client type exists
	var _ = (*ChatClient)(nil)
}
