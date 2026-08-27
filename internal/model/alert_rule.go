package model

import (
	"strings"
	"time"
)

// 告警指标常量。
const (
	MetricRejectRate = "reject_rate"
	MetricLatency    = "latency"
)

// 告警运算符常量。
const (
	OperatorGT = "gt"
	OperatorLT = "lt"
)

// AlertRule 告警规则。
type AlertRule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Resource  string    `json:"resource"`
	Metric    string    `json:"metric"`
	Threshold float64   `json:"threshold"`
	Operator  string    `json:"operator"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 规范化并校验告警规则字段。
func (a *AlertRule) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	a.Resource = strings.TrimSpace(a.Resource)
	if a.Name == "" {
		return NewValidationError("name", "告警规则名称不能为空")
	}
	if a.Metric == "" {
		a.Metric = MetricRejectRate
	}
	if a.Metric != MetricRejectRate && a.Metric != MetricLatency {
		return NewValidationError("metric", "告警指标不合法")
	}
	if a.Operator == "" {
		a.Operator = OperatorGT
	}
	if a.Operator != OperatorGT && a.Operator != OperatorLT {
		return NewValidationError("operator", "告警运算符不合法")
	}
	if a.Threshold < 0 {
		return NewValidationError("threshold", "告警阈值不能为负数")
	}
	if a.Status == "" {
		a.Status = RuleEnabled
	}
	if a.Status != RuleEnabled && a.Status != RuleDisabled {
		return NewValidationError("status", "告警规则状态不合法")
	}
	return nil
}
