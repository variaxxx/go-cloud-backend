package infra_kafka

import "errors"

type nonRetryableError struct {
	err error
}

func NewNonRetryableError(
	err error,
) error {
	if err == nil {
		return nil
	}

	return nonRetryableError{err: err}
}

func (e nonRetryableError) Error() string {
	return e.err.Error()
}

func (e nonRetryableError) Unwrap() error {
	return e.err
}

func IsNonRetryableError(
	err error,
) bool {
	var target nonRetryableError
	return errors.As(err, &target)
}
