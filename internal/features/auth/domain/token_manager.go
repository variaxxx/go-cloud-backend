package auth_domain

type TokenManager interface {
	Parse(tokenString string) (int64, error)
	Issue(userId int64) (string, error)
}
