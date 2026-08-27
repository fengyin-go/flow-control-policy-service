package store

import "errors"

type RetryFailure struct {
	Temporary bool
	Message   string
}

func (e *RetryFailure) Error() string { return e.Message }

type AlertStateRetryRetryState struct {
	steps    []error
	attempts int
	commits  int
}

func NewAlertStateRetryRetryState(steps ...error) *AlertStateRetryRetryState {
	return &AlertStateRetryRetryState{steps: append([]error(nil), steps...)}
}

func (s *AlertStateRetryRetryState) Next() error {
	s.attempts++
	var err error
	if len(s.steps) > 0 {
		err = s.steps[0]
		s.steps = s.steps[1:]
	}
	if err != nil {
		return errors.New(err.Error())
	}
	s.commits++
	return nil
}

func (s *AlertStateRetryRetryState) Attempts() int { return s.attempts }
func (s *AlertStateRetryRetryState) Commits() int  { return s.commits }
