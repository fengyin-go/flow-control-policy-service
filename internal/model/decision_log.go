package model

import (
	"strings"
	"time"
)

// DecisionLog 请求放行决策日志。
type DecisionLog struct {
	ID        string    `json:"id"`
	Resource  string    `json:"resource"`
	Key       string    `json:"key"`
	Allowed   bool      `json:"allowed"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 规范化并校验决策日志字段。
func (d *DecisionLog) Validate() error {
	d.Resource = strings.TrimSpace(d.Resource)
	d.Key = strings.TrimSpace(d.Key)
	d.Reason = strings.TrimSpace(d.Reason)
	if d.Resource == "" {
		return NewValidationError("resource", "资源标识不能为空")
	}
	return nil
}

// DecisionLogFilter 决策日志筛选条件。
type DecisionLogFilter struct {
	Resource string
	Allowed  *bool
}

// Match 判断决策日志是否命中筛选条件。
func (f DecisionLogFilter) Match(d *DecisionLog) bool {
	if f.Resource != "" && d.Resource != f.Resource {
		return false
	}
	if f.Allowed != nil && d.Allowed != *f.Allowed {
		return false
	}
	return true
}
