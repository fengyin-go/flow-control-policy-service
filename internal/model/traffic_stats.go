package model

import (
	"strings"
	"time"
)

// TrafficStats 流量统计。
type TrafficStats struct {
	ID               string    `json:"id"`
	Resource         string    `json:"resource"`
	TotalRequests    int64     `json:"total_requests"`
	AllowedRequests  int64     `json:"allowed_requests"`
	RejectedRequests int64     `json:"rejected_requests"`
	WindowStart      time.Time `json:"window_start"`
	WindowEnd        time.Time `json:"window_end"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Validate 规范化并校验流量统计字段。
func (t *TrafficStats) Validate() error {
	t.Resource = strings.TrimSpace(t.Resource)
	if t.Resource == "" {
		return NewValidationError("resource", "资源标识不能为空")
	}
	if t.TotalRequests < 0 {
		t.TotalRequests = 0
	}
	if t.AllowedRequests < 0 {
		t.AllowedRequests = 0
	}
	if t.RejectedRequests < 0 {
		t.RejectedRequests = 0
	}
	return nil
}
