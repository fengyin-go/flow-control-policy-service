package service

import (
	"sort"
	"time"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/idgen"
)

// CreateTokenBucket 为限流规则创建令牌桶。
func (s *Service) CreateTokenBucket(input model.TokenBucket) (*model.TokenBucket, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetRateLimitRule(input.RuleID); err != nil {
		return nil, model.NewValidationError("rule_id", "限流规则不存在")
	}
	t := &model.TokenBucket{
		ID:         idgen.Hex(),
		RuleID:     input.RuleID,
		Capacity:   input.Capacity,
		RefillRate: input.RefillRate,
		Tokens:     input.Tokens,
		UpdatedAt:  time.Now(),
	}
	if t.Tokens == 0 {
		t.Tokens = float64(t.Capacity)
	}
	if err := s.store.CreateTokenBucket(t); err != nil {
		return nil, err
	}
	return t, nil
}

// ListTokenBuckets 返回全部令牌桶。
func (s *Service) ListTokenBuckets() []*model.TokenBucket {
	list := s.store.ListTokenBuckets()
	sort.Slice(list, func(i, j int) bool { return list[i].RuleID < list[j].RuleID })
	return list
}

// GetTokenBucket 按 ID 获取令牌桶。
func (s *Service) GetTokenBucket(id string) (*model.TokenBucket, error) {
	return s.store.GetTokenBucket(id)
}

// RefillTokenBucket 手动将令牌桶补充到满。
func (s *Service) RefillTokenBucket(id string) (*model.TokenBucket, error) {
	t, err := s.store.GetTokenBucket(id)
	if err != nil {
		return nil, err
	}
	t.Tokens = float64(t.Capacity)
	t.UpdatedAt = time.Now()
	if err := s.store.UpdateTokenBucket(t); err != nil {
		return nil, err
	}
	return t, nil
}

// TakeToken 从令牌桶取出一个令牌。
func (s *Service) TakeToken(id string) (bool, error) {
	t, err := s.store.GetTokenBucket(id)
	if err != nil {
		return false, err
	}
	if t.Tokens < 1 {
		return false, nil
	}
	t.Tokens--
	t.UpdatedAt = time.Now()
	return true, s.store.UpdateTokenBucket(t)
}

// DeleteTokenBucket 删除令牌桶。
func (s *Service) DeleteTokenBucket(id string) error {
	return s.store.DeleteTokenBucket(id)
}
