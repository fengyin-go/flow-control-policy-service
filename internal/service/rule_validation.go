package service

import (
	"strings"

	"flowcontrol/internal/model"
)

// ResourcePattern 校验资源标识的合法性。
func ValidResourcePattern(resource string) bool {
	resource = strings.TrimSpace(resource)
	if resource == "" {
		return false
	}
	// 资源标识需以 / 开头或为纯标识符
	return strings.HasPrefix(resource, "/") || isIdentifier(resource)
}

func isIdentifier(s string) bool {
	for _, r := range s {
		if !(r == '_' || r == '-' || r == '.' || (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return s != ""
}

// ValidateRuleOverlap 检测是否存在资源标识重复的规则。
func (s *Service) ValidateRuleOverlap(resource string, excludeID string) bool {
	for _, r := range s.store.ListRateLimitRules() {
		if r.ID != excludeID && r.Resource == resource {
			return true
		}
	}
	return false
}

// RuleValidationResult 规则校验结果。
type RuleValidationResult struct {
	Valid       bool     `json:"valid"`
	Errors      []string `json:"errors"`
}

// ValidateRuleInput 校验规则输入（资源合法性 + 重叠检测）。
func (s *Service) ValidateRuleInput(input model.RateLimitRule, excludeID string) RuleValidationResult {
	res := RuleValidationResult{Valid: true, Errors: []string{}}
	if !ValidResourcePattern(input.Resource) {
		res.Valid = false
		res.Errors = append(res.Errors, "资源标识不合法")
	}
	if s.ValidateRuleOverlap(input.Resource, excludeID) {
		res.Valid = false
		res.Errors = append(res.Errors, "已存在相同资源标识的规则")
	}
	return res
}
