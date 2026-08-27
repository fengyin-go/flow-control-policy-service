package service

import (
	"sort"
	"time"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/idgen"
)

// CreateRateLimitRule 创建限流规则。
func (s *Service) CreateRateLimitRule(input model.RateLimitRule) (*model.RateLimitRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	r := &model.RateLimitRule{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Resource:  input.Resource,
		Algorithm: input.Algorithm,
		Limit:     input.Limit,
		WindowSec: input.WindowSec,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateRateLimitRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

// GetRateLimitRule 按 ID 获取限流规则。
func (s *Service) GetRateLimitRule(id string) (*model.RateLimitRule, error) {
	return s.store.GetRateLimitRule(id)
}

// ListRateLimitRules 按筛选条件分页查询限流规则。
func (s *Service) ListRateLimitRules(filter model.RuleFilter, page, size int) ([]*model.RateLimitRule, int, error) {
	all := s.store.ListRateLimitRules()
	matched := make([]*model.RateLimitRule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.Before(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RateLimitRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateRateLimitRule 更新限流规则。
func (s *Service) UpdateRateLimitRule(id string, input model.RateLimitRule) (*model.RateLimitRule, error) {
	existing, err := s.store.GetRateLimitRule(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Resource = input.Resource
	existing.Algorithm = input.Algorithm
	existing.Limit = input.Limit
	existing.WindowSec = input.WindowSec
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateRateLimitRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// ToggleRateLimitRule 启用或停用限流规则。
func (s *Service) ToggleRateLimitRule(id string, enabled bool) (*model.RateLimitRule, error) {
	r, err := s.store.GetRateLimitRule(id)
	if err != nil {
		return nil, err
	}
	if enabled {
		r.Status = model.RuleEnabled
	} else {
		r.Status = model.RuleDisabled
	}
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRateLimitRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteRateLimitRule 删除限流规则及其令牌桶。
func (s *Service) DeleteRateLimitRule(id string) error {
	if _, err := s.store.GetRateLimitRule(id); err != nil {
		return err
	}
	if bucket, err := s.store.GetTokenBucketByRule(id); err == nil {
		_ = s.store.DeleteTokenBucket(bucket.ID)
	}
	return s.store.DeleteRateLimitRule(id)
}

// AllowRequest 按资源判定请求是否被限流。
func (s *Service) AllowRequest(resource, key string) (bool, string) {
	rule, err := s.store.GetRateLimitRuleByResource(resource)
	if err != nil || rule.Status != model.RuleEnabled {
		return true, "no_rule"
	}
	if s.allowByRule(rule, key) {
		return true, "allowed"
	}
	return false, "rate_limited"
}
