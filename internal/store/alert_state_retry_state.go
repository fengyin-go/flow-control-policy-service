package store

// RetryFailure 表示一次告警状态写入失败，通过 Temporary 区分永久与临时错误。
// Temporary=true 表示瞬时繁忙（如锁竞争、瞬时不可用），可重试；false 表示永久拒绝
// （如校验失败、状态冲突），重试无意义且不应再进入下一次写入。
type RetryFailure struct {
	Temporary bool
	Message   string
}

func (e *RetryFailure) Error() string { return e.Message }

// IsTemporary 返回 err 是否为临时错误；非 *RetryFailure 的错误视为不可重试。
func IsTemporary(err error) bool {
	if err == nil {
		return false
	}
	f, ok := err.(*RetryFailure)
	return ok && f.Temporary
}

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
		// 保留原始错误类型（含 Temporary 标志），供上层按永久/临时区分重试策略；
		// 不可用 errors.New(err.Error()) 抹平类型，否则永久拒绝也会进入下一次写入。
		if _, ok := err.(*RetryFailure); ok {
			return err
		}
		return &RetryFailure{Temporary: false, Message: err.Error()}
	}
	s.commits++
	return nil
}

func (s *AlertStateRetryRetryState) Attempts() int { return s.attempts }
func (s *AlertStateRetryRetryState) Commits() int  { return s.commits }
