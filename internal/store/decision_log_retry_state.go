package store

import "errors"

type RetryFailure struct {
	Temporary bool
	Message   string
}

func (e *RetryFailure) Error() string { return e.Message }

// IsTemporary 报告 err 是否为可重试的临时错误。
// 只有显式标记为 Temporary 的 RetryFailure 才视为可重试；其余错误一律按永久错误处理。
func IsTemporary(err error) bool {
	var f *RetryFailure
	if errors.As(err, &f) {
		return f.Temporary
	}
	return false
}

type DecisionLogRetryRetryState struct {
	steps    []error
	attempts int
	commits  int
	lastErr  error
}

func NewDecisionLogRetryRetryState(steps ...error) *DecisionLogRetryRetryState {
	return &DecisionLogRetryRetryState{steps: append([]error(nil), steps...)}
}

// Next 尝试一次写入，结果分三类：
//   - nil：写入成功（commits 自增）。
//   - 临时错误：可重试，不计入 commits，消耗该步骤后供下一次重试。
//   - 永久错误：立即停止，不消耗后续步骤、不计入 commits。
//
// 步骤耗尽时返回最近一次临时错误（重试耗尽），而非当作成功，
// 以免"仅出现临时错误"的场景被误判为已提交。
func (s *DecisionLogRetryRetryState) Next() error {
	s.attempts++
	if len(s.steps) == 0 {
		// 无更多步骤可重试，返回最近一次错误（重试耗尽）。
		return s.lastErr
	}
	err := s.steps[0]
	s.steps = s.steps[1:]
	if err == nil {
		// 该步骤表示一次成功提交。
		s.commits++
		s.lastErr = nil
		return nil
	}
	if !IsTemporary(err) {
		// 永久错误：停止消耗后续步骤，也不计入 commits。
		s.lastErr = err
		s.steps = nil
		return err
	}
	// 临时错误：保留为最近错误，供重试耗尽时返回。
	s.lastErr = err
	return err
}

// Last 返回最近一次错误（成功后为 nil），供重试耗尽场景查询。
func (s *DecisionLogRetryRetryState) Last() error { return s.lastErr }

func (s *DecisionLogRetryRetryState) Attempts() int { return s.attempts }
func (s *DecisionLogRetryRetryState) Commits() int  { return s.commits }
