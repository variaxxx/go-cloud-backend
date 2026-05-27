package auth_password

import "testing"

func TestHasherHashAndCompare(t *testing.T) {
	hasher := NewHasher()

	hash, err := hasher.Hash("secret")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if hash == "" || hash == "secret" {
		t.Fatalf("Hash() = %q, want non-empty bcrypt hash", hash)
	}
	if err := hasher.Compare(hash, "secret"); err != nil {
		t.Fatalf("Compare() error = %v", err)
	}
	if err := hasher.Compare(hash, "wrong"); err == nil {
		t.Fatal("Compare() error = nil, want mismatch")
	}
}
