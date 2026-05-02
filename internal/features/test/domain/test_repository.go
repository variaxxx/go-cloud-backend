package test_domain

import "context"

type TestRepository interface {
	WriteStr(
		ctx context.Context,
		str string,
	) error
}
