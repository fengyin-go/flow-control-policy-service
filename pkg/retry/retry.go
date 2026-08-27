// Package retry 提供带指数退避的通用重试能力。
package retry

import (
	"context"
	"time"

	"flowcontrol/pkg/backoff"
)

// Do 执行 fn，失败时按退避策略重试，直至成功或超出 attempts 次。
func Do(ctx context.Context, attempts int, fn func() error) error {
	b := backoff.New()
	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(b.Next()):
		}
	}
	return lastErr
}
