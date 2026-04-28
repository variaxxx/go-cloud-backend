package auth_password

import "golang.org/x/crypto/bcrypt"

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Hash(
	password string,
) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (h *Hasher) Compare(
	hash string,
	password string,
) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
