package service

import (
	"sort"
	"time"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/idgen"
)

// CreateAlertRule 创建告警规则。
func (s *Service) CreateAlertRule(input model.AlertRule) (*model.AlertRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a := &model.AlertRule{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Resource:  input.Resource,
		Metric:    input.Metric,
		Threshold: input.Threshold,
		Operator:  input.Operator,
		Status:    input.Status,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateAlertRule(a); err != nil {
		return nil, err
	}
	return a, nil
}

// GetAlertRule 按 ID 获取告警规则。
func (s *Service) GetAlertRule(id string) (*model.AlertRule, error) {
	return s.store.GetAlertRule(id)
}

// ListAlertRules 返回全部告警规则。
func (s *Service) ListAlertRules() []*model.AlertRule {
	list := s.store.ListAlertRules()
	sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
	return list
}

// UpdateAlertRule 更新告警规则。
func (s *Service) UpdateAlertRule(id string, input model.AlertRule) (*model.AlertRule, error) {
	existing, err := s.store.GetAlertRule(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Resource = input.Resource
	existing.Metric = input.Metric
	existing.Threshold = input.Threshold
	existing.Operator = input.Operator
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAlertRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteAlertRule 删除告警规则。
func (s *Service) DeleteAlertRule(id string) error {
	return s.store.DeleteAlertRule(id)
}

// AlertHit 告警命中结果。
type AlertHit struct {
	Rule      *model.AlertRule `json:"rule"`
	MetricValue float64        `json:"metric_value"`
}

// EvaluateAlerts 评估全部启用的告警规则，返回命中的规则。
func (s *Service) EvaluateAlerts() []AlertHit {
	hits := make([]AlertHit, 0)
	for _, rule := range s.store.ListAlertRules() {
		if rule.Status != model.RuleEnabled {
			continue
		}
		var metricValue float64
		if stats, err := s.store.GetTrafficStatsByResource(rule.Resource); err == nil {
			if rule.Metric == model.MetricRejectRate {
				if stats.TotalRequests > 0 {
					metricValue = float64(stats.RejectedRequests) / float64(stats.TotalRequests)
				}
			}
		}
		hit := false
		if rule.Operator == model.OperatorGT && metricValue > rule.Threshold {
			hit = true
		}
		if rule.Operator == model.OperatorLT && metricValue < rule.Threshold {
			hit = true
		}
		if hit {
			hits = append(hits, AlertHit{Rule: rule, MetricValue: metricValue})
		}
	}
	return hits
}
