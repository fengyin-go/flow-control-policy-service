package store

import "errors"

// RetryFailure 表示一次流量统计写入的失败结果。
// Temporary 为 true 时属于临时故障，可重试；为 false 时属于永久拒绝，必须停止重试。
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
		// 保留原始错误类型（含 RetryFailure.Temporary），以便调用方区分临时故障与
		// 永久拒绝；切勿用 errors.New(err.Error()) 抹掉类型，否则永久拒绝会被误判为可重试。
		return err
	}
	s.commits++
	return nil
}

func (s *TrafficStatsRetryRetryState) Attempts() int { return s.attempts }
func (s *TrafficStatsRetryRetryState) Commits() int  { return s.commits }

// IsPermanentFailure 报告 err 是否为不可重试的永久拒绝。
func IsPermanentFailure(err error) bool {
	var f *RetryFailure
	if !errors.As(err, &f) {
		return false
	}
	return !f.Temporary
}
