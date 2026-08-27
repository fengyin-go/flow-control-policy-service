package model

import (
	"strings"
	"time"
)

// 熔断器状态常量。
const (
	BreakerClosed   = "closed"
	BreakerOpen     = "open"
	BreakerHalfOpen = "half_open"
)

// breakerTransitions 定义熔断器合法状态流转。
var breakerTransitions = map[string]map[string]bool{
	BreakerClosed:   {BreakerOpen: true},
	BreakerOpen:     {BreakerHalfOpen: true},
	BreakerHalfOpen: {BreakerClosed: true, BreakerOpen: true},
}

// CanBreakerTransition 判断熔断器能否从 from 流转到 to。
func CanBreakerTransition(from, to string) bool {
	if m, ok := breakerTransitions[from]; ok {
		return m[to]
	}
	return false
}

// CircuitBreaker 熔断器。
type CircuitBreaker struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Resource         string    `json:"resource"`
	FailureThreshold int       `json:"failure_threshold"`
	SuccessThreshold int       `json:"success_threshold"`
	TimeoutMs        int       `json:"timeout_ms"`
	OpenDurationSec  int       `json:"open_duration_sec"`
	State            string    `json:"state"`
	FailCount        int       `json:"fail_count"`
	SuccessCount     int       `json:"success_count"`
	LastStateChange  time.Time `json:"last_state_change"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Validate 规范化并校验熔断器字段。
func (c *CircuitBreaker) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Resource = strings.TrimSpace(c.Resource)
	if c.Name == "" {
		return NewValidationError("name", "熔断器名称不能为空")
	}
	if c.Resource == "" {
		return NewValidationError("resource", "资源标识不能为空")
	}
	if c.FailureThreshold <= 0 {
		c.FailureThreshold = 5
	}
	if c.SuccessThreshold <= 0 {
		c.SuccessThreshold = 2
	}
	if c.TimeoutMs <= 0 {
		c.TimeoutMs = 1000
	}
	if c.OpenDurationSec <= 0 {
		c.OpenDurationSec = 30
	}
	if c.State == "" {
		c.State = BreakerClosed
	}
	if c.State != BreakerClosed && c.State != BreakerOpen && c.State != BreakerHalfOpen {
		return NewValidationError("state", "熔断器状态不合法")
	}
	if c.FailCount < 0 {
		c.FailCount = 0
	}
	if c.SuccessCount < 0 {
		c.SuccessCount = 0
	}
	return nil
}
