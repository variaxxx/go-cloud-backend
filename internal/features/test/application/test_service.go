package test_app

import (
	test_domain "cloud/internal/features/test/domain"
	"context"
)

type testService struct {
	repository test_domain.TestRepository
}

func NewTestService(
	repository test_domain.TestRepository,
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
