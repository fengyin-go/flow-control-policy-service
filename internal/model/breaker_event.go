package model

import (
	"strings"
	"time"
)

// 熔断事件类型常量。
const (
	BreakerEventOpened    = "opened"
	BreakerEventClosed    = "closed"
	BreakerEventHalfOpen  = "half_open"
	BreakerEventRejected  = "rejected"
)

// BreakerEvent 熔断事件。
type BreakerEvent struct {
	ID        string    `json:"id"`
	BreakerID string    `json:"breaker_id"`
	EventType string    `json:"event_type"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 规范化并校验熔断事件字段。
func (e *BreakerEvent) Validate() error {
	e.BreakerID = strings.TrimSpace(e.BreakerID)
	e.EventType = strings.TrimSpace(e.EventType)
	e.Detail = strings.TrimSpace(e.Detail)
	if e.EventType == "" {
		return NewValidationError("event_type", "事件类型不能为空")
	}
	return nil
}
