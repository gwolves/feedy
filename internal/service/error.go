package service

type ExpectedError interface {
	Reason() string
}

type ErrorCode int

const (
	NotFound ErrorCode = iota
	InvalidFeed
)

func WithReason(err error, reason string) error {
	return expectedError{
		error:  err,
		reason: reason,
	}
}

type expectedError struct {
	reason string
	error
}

func (e expectedError) Reason() string {
	return e.reason
}
