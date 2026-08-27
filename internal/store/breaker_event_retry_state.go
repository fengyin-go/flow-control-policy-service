package store

type RetryFailure struct {
	Temporary bool
	Message   string
}

func (e *RetryFailure) Error() string { return e.Message }

type BreakerEventRetryRetryState struct {
	steps    []error
	attempts int
	commits  int
}

func NewBreakerEventRetryRetryState(steps ...error) *BreakerEventRetryRetryState {
	return &BreakerEventRetryRetryState{steps: append([]error(nil), steps...)}
}

// Next 返回下一条预置错误，保留原始错误类别（如 *RetryFailure 的 Temporary 标记），
// 供上层据此判断是否可重试。步骤耗尽或遇到 nil 步骤时计为一次提交。
func (s *BreakerEventRetryRetryState) Next() error {
	s.attempts++
	var err error
	if len(s.steps) > 0 {
		err = s.steps[0]
		s.steps = s.steps[1:]
	}
	if err != nil {
		return err
	}
	s.commits++
	return nil
}

func (s *BreakerEventRetryRetryState) Attempts() int { return s.attempts }
func (s *BreakerEventRetryRetryState) Commits() int  { return s.commits }
