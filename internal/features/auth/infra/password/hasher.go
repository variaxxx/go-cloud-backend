package auth_password

import "golang.org/x/crypto/bcrypt"

type hasher struct{}

func NewHasher() *hasher {
	return &hasher{}
}

func (h *hasher) Hash(
	password string,
) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (h *hasher) Compare(
	hash string,
	password string,
) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
