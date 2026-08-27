package store

import (
	"sync"

	"flowcontrol/internal/model"
)

// MemoryStore 基于内存的 Store 实现，线程安全。
type MemoryStore struct {
	mu             sync.RWMutex
	rules          map[string]*model.RateLimitRule
	tokenBuckets   map[string]*model.TokenBucket
	breakers       map[string]*model.CircuitBreaker
	breakerEvents  map[string]*model.BreakerEvent
	quotas         map[string]*model.Quota
	trafficStats   map[string]*model.TrafficStats
	alertRules     map[string]*model.AlertRule
	decisionLogs   map[string]*model.DecisionLog
}

// NewMemoryStore 构造空的内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rules:         make(map[string]*model.RateLimitRule),
		tokenBuckets:  make(map[string]*model.TokenBucket),
		breakers:      make(map[string]*model.CircuitBreaker),
		breakerEvents: make(map[string]*model.BreakerEvent),
		quotas:        make(map[string]*model.Quota),
		trafficStats:  make(map[string]*model.TrafficStats),
		alertRules:    make(map[string]*model.AlertRule),
		decisionLogs:  make(map[string]*model.DecisionLog),
	}
}

var _ Store = (*MemoryStore)(nil)
