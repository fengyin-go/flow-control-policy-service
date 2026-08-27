package service

// DecisionResult 请求放行决策结果。
type DecisionResult struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
	Rule    string `json:"rule,omitempty"`
	Breaker string `json:"breaker,omitempty"`
}

// Decide 综合熔断器与限流规则给出放行决策，并更新流量统计与决策日志。
func (s *Service) Decide(resource, key string) DecisionResult {
	if allowed, _ := s.CheckBreaker(resource); !allowed {
		s.recordTraffic(resource, false)
		s.recordDecision(resource, key, false, "circuit_open")
		return DecisionResult{Allowed: false, Reason: "circuit_open", Breaker: resource}
	}
	allowed, reason := s.AllowRequest(resource, key)
	s.recordTraffic(resource, allowed)
	s.recordDecision(resource, key, allowed, reason)
	return DecisionResult{Allowed: allowed, Reason: reason, Rule: resource}
}
