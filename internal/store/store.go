// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"flowcontrol/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// 限流规则
	CreateRateLimitRule(r *model.RateLimitRule) error
	GetRateLimitRule(id string) (*model.RateLimitRule, error)
	GetRateLimitRuleByResource(resource string) (*model.RateLimitRule, error)
	ListRateLimitRules() []*model.RateLimitRule
	UpdateRateLimitRule(r *model.RateLimitRule) error
	DeleteRateLimitRule(id string) error

	// 令牌桶
	CreateTokenBucket(t *model.TokenBucket) error
	GetTokenBucket(id string) (*model.TokenBucket, error)
	GetTokenBucketByRule(ruleID string) (*model.TokenBucket, error)
	ListTokenBuckets() []*model.TokenBucket
	UpdateTokenBucket(t *model.TokenBucket) error
	DeleteTokenBucket(id string) error

	// 熔断器
	CreateCircuitBreaker(c *model.CircuitBreaker) error
	GetCircuitBreaker(id string) (*model.CircuitBreaker, error)
	GetCircuitBreakerByResource(resource string) (*model.CircuitBreaker, error)
	ListCircuitBreakers() []*model.CircuitBreaker
	UpdateCircuitBreaker(c *model.CircuitBreaker) error
	DeleteCircuitBreaker(id string) error

	// 熔断事件
	CreateBreakerEvent(e *model.BreakerEvent) error
	ListBreakerEvents() []*model.BreakerEvent

	// 配额
	CreateQuota(q *model.Quota) error
	GetQuota(id string) (*model.Quota, error)
	ListQuotas() []*model.Quota
	UpdateQuota(q *model.Quota) error
	DeleteQuota(id string) error

	// 流量统计
	CreateTrafficStats(t *model.TrafficStats) error
	GetTrafficStats(id string) (*model.TrafficStats, error)
	GetTrafficStatsByResource(resource string) (*model.TrafficStats, error)
	ListTrafficStats() []*model.TrafficStats
	UpdateTrafficStats(t *model.TrafficStats) error

	// 告警规则
	CreateAlertRule(a *model.AlertRule) error
	GetAlertRule(id string) (*model.AlertRule, error)
	ListAlertRules() []*model.AlertRule
	UpdateAlertRule(a *model.AlertRule) error
	DeleteAlertRule(id string) error

	// 决策日志
	CreateDecisionLog(d *model.DecisionLog) error
	ListDecisionLogs() []*model.DecisionLog
}
