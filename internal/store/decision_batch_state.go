package store

import "errors"

// ErrStreamInterrupted 是放行决策批次遇错中断时返回的哨兵错误。
// 导出为变量，便于调用方用 errors.Is 精确识别批次中断。
var ErrStreamInterrupted = errors.New("stream interrupted")

type DecisionBatchStream struct{}

func NewDecisionBatchStream() *DecisionBatchStream { return &DecisionBatchStream{} }

func (s *DecisionBatchStream) Start(items []string, failAt int) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1)
	go func() {

		// 无论是正常结束还是遇错中断，都必须关闭 out，否则消费者的
		// `for range out` 会因永远收不到关闭信号而死锁，整个批次挂死。
		defer close(out)
		defer close(errs)
		for index, item := range items {
			if index == failAt {
				errs <- ErrStreamInterrupted
				return
			}
			out <- item
		}
	}()
	return out, errs
}
