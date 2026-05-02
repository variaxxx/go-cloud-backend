package test_app

import "context"

type TestUseCase interface {
	GetHello() string

	WriteStr(
		ctx context.Context,
		str string,
	) error
}
