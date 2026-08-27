package service

import (
	"sort"
	"time"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/idgen"
)

// CreateQuota 创建配额。
func (s *Service) CreateQuota(input model.Quota) (*model.Quota, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetRateLimitRule(input.RuleID); err != nil {
		return nil, model.NewValidationError("rule_id", "限流规则不存在")
	}
	q := &model.Quota{
		ID:             idgen.Hex(),
		RuleID:         input.RuleID,
		Dimension:      input.Dimension,
		DimensionValue: input.DimensionValue,
		Allowed:        input.Allowed,
		Used:           0,
		ResetAt:        time.Now().Add(time.Minute),
	}
	if err := s.store.CreateQuota(q); err != nil {
		return nil, err
	}
	return q, nil
}

// ListQuotas 返回配额（可按规则筛选）。
func (s *Service) ListQuotas(ruleID string) []*model.Quota {
	list := make([]*model.Quota, 0)
	for _, q := range s.store.ListQuotas() {
		if ruleID == "" || q.RuleID == ruleID {
			list = append(list, q)
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].DimensionValue < list[j].DimensionValue })
	return list
}

// ConsumeQuota 消耗配额，返回是否还有剩余。
func (s *Service) ConsumeQuota(ruleID, dimension, value string) (bool, error) {
	for _, q := range s.store.ListQuotas() {
		if q.RuleID == ruleID && q.Dimension == dimension && q.DimensionValue == value {
			if q.Used >= q.Allowed {
				return false, nil
			}
			q.Used++
			return true, s.store.UpdateQuota(q)
		}
	}
	return false, model.NewValidationError("quota", "配额不存在")
}

// DeleteQuota 删除配额。
func (s *Service) DeleteQuota(id string) error {
	return s.store.DeleteQuota(id)
}
