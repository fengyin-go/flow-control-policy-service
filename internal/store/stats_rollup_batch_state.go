package store

import "errors"

type StatsRollupBatchStream struct{}

func NewStatsRollupBatchStream() *StatsRollupBatchStream { return &StatsRollupBatchStream{} }

func (s *StatsRollupBatchStream) Start(items []string, failAt int) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1)
	go func() {
		// 关闭 out，使消费者在流出错提前返回或正常结束后能退出 range，
		// 否则归并方已收到部分结果却会因 out 永不关闭而卡死。
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
