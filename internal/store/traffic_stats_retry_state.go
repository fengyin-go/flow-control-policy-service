package store

import "errors"

type RetryFailure struct {
	Temporary bool
	Message   string
}

func (e *RetryFailure) Error() string { return e.Message }

type TrafficStatsRetryRetryState struct {
	steps    []error
	attempts int
	commits  int
}

func NewTrafficStatsRetryRetryState(steps ...error) *TrafficStatsRetryRetryState {
	return &TrafficStatsRetryRetryState{steps: append([]error(nil), steps...)}
}

func (s *TrafficStatsRetryRetryState) Next() error {
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

func (s *TrafficStatsRetryRetryState) Attempts() int { return s.attempts }
func (s *TrafficStatsRetryRetryState) Commits() int  { return s.commits }
