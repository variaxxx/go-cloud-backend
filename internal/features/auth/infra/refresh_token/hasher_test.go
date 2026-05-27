package auth_refresh_token

import "testing"

func TestHasherHash(t *testing.T) {
	hasher := NewHasher()

	got := hasher.Hash("refresh-token")
	again := hasher.Hash("refresh-token")
	other := hasher.Hash("other-token")

	if got == "" {
		t.Fatal("Hash() is empty")
	}
	if got != again {
		t.Fatalf("Hash() is not deterministic: %q != %q", got, again)
	}
	if got == other {
		t.Fatalf("Hash() collision for different test inputs: %q", got)
	}
	if len(got) != 64 {
		t.Fatalf("Hash() length = %d, want 64 hex chars", len(got))
	}
}
