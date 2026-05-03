package test_app

import "context"

type TestRepository interface {
	WriteStr(
		ctx context.Context,
		str string,
	) error
}
