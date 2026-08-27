package store

import "errors"

type QuotaConsumeBatchStream struct{}

func NewQuotaConsumeBatchStream() *QuotaConsumeBatchStream { return &QuotaConsumeBatchStream{} }

func (s *QuotaConsumeBatchStream) Start(items []string, failAt int) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1)
	go func() {
		// 失败或正常结束时都要关闭结果流，否则消费方的 range 会永久阻塞，
		// 无法走到读取 errs 并返回部分结果那一步。
		defer close(out)
		defer close(errs)
		for index, item := range items {
			if index == failAt {
				errs <- errors.New("stream interrupted")
				return
			}
			out <- item
		}
	}()
	return out, errs
}
