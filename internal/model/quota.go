package model

import (
	"strings"
	"time"
)

// 配额维度常量。
const (
	DimensionIP   = "ip"
	DimensionUser = "user"
	DimensionAPI  = "api"
)

// Quota 配额（按维度计数）。
type Quota struct {
	ID             string    `json:"id"`
	RuleID         string    `json:"rule_id"`
	Dimension      string    `json:"dimension"`
	DimensionValue string    `json:"dimension_value"`
	Allowed        int       `json:"allowed"`
	Used           int       `json:"used"`
	ResetAt        time.Time `json:"reset_at"`
}

// Validate 规范化并校验配额字段。
func (q *Quota) Validate() error {
	q.RuleID = strings.TrimSpace(q.RuleID)
	q.Dimension = strings.TrimSpace(q.Dimension)
	q.DimensionValue = strings.TrimSpace(q.DimensionValue)
	if q.RuleID == "" {
		return NewValidationError("rule_id", "限流规则 ID 不能为空")
	}
	if q.Dimension == "" {
		q.Dimension = DimensionIP
	}
	if q.DimensionValue == "" {
		return NewValidationError("dimension_value", "配额维度值不能为空")
	}
	if q.Allowed < 0 {
		q.Allowed = 0
	}
	if q.Used < 0 {
		q.Used = 0
	}
	return nil
}
