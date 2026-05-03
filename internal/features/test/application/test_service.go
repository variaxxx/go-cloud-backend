package test_app

import (
	"context"
)

type testService struct {
	repository TestRepository
}

func NewTestService(
	repository TestRepository,
) *testService {
	return &testService{
		repository: repository,
	}
}

func (s *testService) GetHello() string {
	return "Hello!"
}

func (s *testService) WriteStr(
	ctx context.Context,
	str string,
) error {
	return s.repository.WriteStr(ctx, str)
}
