package service

import (
	"flowcontrol/internal/model"
	"flowcontrol/pkg/idgen"
)

// RulePreset 限流规则预设。
type RulePreset struct {
	Name      string `json:"name"`
	Resource  string `json:"resource"`
	Algorithm string `json:"algorithm"`
	Limit     int    `json:"limit"`
	WindowSec int    `json:"window_sec"`
}

// DefaultRulePresets 返回常用限流规则预设。
func DefaultRulePresets() []RulePreset {
	return []RulePreset{
		{Name: "全局接口限流", Resource: "/api", Algorithm: model.AlgorithmFixedWindow, Limit: 100, WindowSec: 1},
		{Name: "登录接口限流", Resource: "/api/login", Algorithm: model.AlgorithmSlidingWindow, Limit: 10, WindowSec: 60},
		{Name: "列表查询限流", Resource: "/api/list", Algorithm: model.AlgorithmTokenBucket, Limit: 50, WindowSec: 1},
		{Name: "上传接口限流", Resource: "/api/upload", Algorithm: model.AlgorithmSlidingWindow, Limit: 5, WindowSec: 60},
	}
}

// SeedDefaultRules 将预设规则写入存储（幂等，已存在则跳过），返回创建数量。
func (s *Service) SeedDefaultRules() int {
	created := 0
	for _, p := range DefaultRulePresets() {
		rule := &model.RateLimitRule{
			ID:        idgen.Hex(),
			Name:      p.Name,
			Resource:  p.Resource,
			Algorithm: p.Algorithm,
			Limit:     p.Limit,
			WindowSec: p.WindowSec,
			Status:    model.RuleEnabled,
		}
		if err := s.store.CreateRateLimitRule(rule); err == nil {
			created++
		}
	}
	return created
}
