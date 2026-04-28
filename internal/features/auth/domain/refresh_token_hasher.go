package auth_domain

type RefreshTokenHasher interface {
	Hash(token string) string
}
