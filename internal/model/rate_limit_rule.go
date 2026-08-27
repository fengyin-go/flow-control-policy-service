package model

import (
	"strings"
	"time"
)

// 限流算法常量。
const (
	AlgorithmFixedWindow   = "fixed_window"
	AlgorithmSlidingWindow = "sliding_window"
	AlgorithmTokenBucket   = "token_bucket"
)

// 限流规则状态常量。
const (
	RuleEnabled  = "enabled"
	RuleDisabled = "disabled"
)

// RateLimitRule 限流规则。
type RateLimitRule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Resource  string    `json:"resource"`
	Algorithm string    `json:"algorithm"`
	Limit     int       `json:"limit"`
	WindowSec int       `json:"window_sec"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 规范化并校验限流规则字段。
func (r *RateLimitRule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Resource = strings.TrimSpace(r.Resource)
	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if r.Resource == "" {
		return NewValidationError("resource", "资源标识不能为空")
	}
	if r.Algorithm == "" {
		r.Algorithm = AlgorithmFixedWindow
	}
	switch r.Algorithm {
	case AlgorithmFixedWindow, AlgorithmSlidingWindow, AlgorithmTokenBucket:
	default:
		return NewValidationError("algorithm", "限流算法不合法")
	}
	if r.Limit <= 0 {
		return NewValidationError("limit", "限流阈值必须大于 0")
	}
	if r.WindowSec <= 0 {
		r.WindowSec = 60
	}
	if r.Status == "" {
		r.Status = RuleEnabled
	}
	if r.Status != RuleEnabled && r.Status != RuleDisabled {
		return NewValidationError("status", "规则状态不合法")
	}
	return nil
}

// RuleFilter 规则筛选条件。
type RuleFilter struct {
	Resource  string
	Algorithm string
	Status    string
	Keyword   string
}

// Match 判断规则是否命中筛选条件。
func (f RuleFilter) Match(r *RateLimitRule) bool {
	if f.Resource != "" && r.Resource != f.Resource {
		return false
	}
	if f.Algorithm != "" && r.Algorithm != f.Algorithm {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) {
			return false
		}
	}
	return true
}
