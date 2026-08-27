package model

import (
	"strings"
	"time"
)

// TokenBucket 令牌桶配置（token_bucket 算法的运行时快照）。
type TokenBucket struct {
	ID         string    `json:"id"`
	RuleID     string    `json:"rule_id"`
	Capacity   int       `json:"capacity"`
	RefillRate int       `json:"refill_rate"`
	Tokens     float64   `json:"tokens"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate 规范化并校验令牌桶字段。
func (t *TokenBucket) Validate() error {
	t.RuleID = strings.TrimSpace(t.RuleID)
	if t.RuleID == "" {
		return NewValidationError("rule_id", "限流规则 ID 不能为空")
	}
	if t.Capacity <= 0 {
		return NewValidationError("capacity", "令牌桶容量必须大于 0")
	}
	if t.RefillRate <= 0 {
		return NewValidationError("refill_rate", "令牌补充速率必须大于 0")
	}
	if t.Tokens < 0 {
		t.Tokens = 0
	}
	if t.Tokens > float64(t.Capacity) {
		t.Tokens = float64(t.Capacity)
	}
	return nil
}
