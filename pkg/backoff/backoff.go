// Package backoff 提供带抖动的指数退避策略。
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Exponential 指数退避策略。
type Exponential struct {
	Base     time.Duration
	Max      time.Duration
	Factor   float64
	Jitter   float64
	attempts int
}

// New 构造默认退避策略：基数 200ms，上限 10s，因子 2，抖动 0.2。
func New() *Exponential {
	return &Exponential{
		Base:   200 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2,
		Jitter: 0.2,
	}
}

// Next 返回下一次退避时长，并累加尝试次数。
func (e *Exponential) Next() time.Duration {
	if e.Base <= 0 {
		e.Base = 200 * time.Millisecond
	}
	if e.Factor <= 1 {
		e.Factor = 2
	}
	d := float64(e.Base) * math.Pow(e.Factor, float64(e.attempts))
	if e.Max > 0 && d > float64(e.Max) {
		d = float64(e.Max)
	}
	if e.Jitter > 0 {
		jitter := d * e.Jitter
		d = d - jitter + rand.Float64()*2*jitter
	}
	e.attempts++
	return time.Duration(d)
}

// Reset 重置尝试次数。
func (e *Exponential) Reset() {
	e.attempts = 0
}

// Attempts 返回当前尝试次数。
func (e *Exponential) Attempts() int {
	return e.attempts
}
