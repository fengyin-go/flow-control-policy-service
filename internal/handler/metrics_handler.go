package handler

import (
	"net/http"
	"strconv"

	"flowcontrol/internal/service"
	"flowcontrol/pkg/httpx"
)

func (s *Server) registerMetricsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/limiter/metrics", s.limiterMetrics)
	mux.HandleFunc("POST /api/limiter/reset", s.resetLimiter)
	mux.HandleFunc("GET /api/breakers/top-active", s.mostActiveBreakers)
	mux.HandleFunc("GET /api/breakers/health-summary", s.breakerHealthSummary)
	mux.HandleFunc("GET /api/breakers/recoveries", s.breakerRecoveries)
	mux.HandleFunc("GET /api/breakers/{id}/metrics", s.breakerMetrics)
	mux.HandleFunc("GET /api/breakers/{id}/timeline", s.breakerTimeline)
	mux.HandleFunc("GET /api/breakers/{id}/event-analysis", s.breakerEventAnalysis)
	mux.HandleFunc("GET /api/health-report", s.healthReport)
	mux.HandleFunc("GET /api/reports/rule-distribution", s.ruleDistribution)
	mux.HandleFunc("GET /api/reports/rule-status", s.ruleStatusReport)
	mux.HandleFunc("GET /api/reports/top-rejected", s.topRejectedResources)
	mux.HandleFunc("GET /api/reports/breaker-events", s.breakerEventDistribution)
	mux.HandleFunc("GET /api/reports/quota-usage", s.topQuotaUsage)
	mux.HandleFunc("GET /api/reports/quota-report", s.quotaReport)
	mux.HandleFunc("GET /api/reports/thresholds", s.adviseThresholds)
	mux.HandleFunc("GET /api/reports/resource-statuses", s.resourceStatuses)
	mux.HandleFunc("GET /api/reports/resource-groups", s.resourceGroups)
	mux.HandleFunc("GET /api/reports/decision-trends", s.decisionTrends)
	mux.HandleFunc("GET /api/reports/top-rejected-keys", s.topRejectedKeys)
	mux.HandleFunc("GET /api/reports/reason-distribution", s.reasonDistribution)
	mux.HandleFunc("GET /api/reports/alert-analysis", s.alertAnalysis)
	mux.HandleFunc("GET /api/reports/quota-analysis", s.quotaAnalysis)
	mux.HandleFunc("GET /api/reports/rule-recommendations", s.ruleRecommendations)
	mux.HandleFunc("GET /api/reports/config", s.configReport)
	mux.HandleFunc("GET /api/reports/optimization", s.optimizationReport)
	mux.HandleFunc("GET /api/reports/rule-complexity", s.ruleComplexityReport)
	mux.HandleFunc("GET /api/reports/breaker-tuning", s.breakerTuningSuggestions)
	mux.HandleFunc("GET /api/reports/breaker-effectiveness", s.breakerEffectiveness)
	mux.HandleFunc("GET /api/reports/traffic-forecast", s.trafficForecasts)
	mux.HandleFunc("GET /api/reports/rule-stats", s.ruleStatsReport)
	mux.HandleFunc("GET /api/reports/disabled-rules", s.disabledRuleResources)
	mux.HandleFunc("GET /api/reports/inventory", s.resourceInventory)
	mux.HandleFunc("GET /api/summary", s.systemSummary)
	mux.HandleFunc("GET /api/rule-presets", s.rulePresets)
	mux.HandleFunc("POST /api/rule-presets/seed", s.seedDefaultRules)
	mux.HandleFunc("POST /api/simulate", s.simulateTraffic)
}

func (s *Server) limiterMetrics(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.LimiterRuntimeMetrics())
}

func (s *Server) resetLimiter(w http.ResponseWriter, r *http.Request) {
	s.svc.ResetLimiter()
	httpx.OK(w, map[string]string{"status": "reset"})
}

func (s *Server) mostActiveBreakers(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	httpx.OK(w, s.svc.MostActiveBreakers(n))
}

func (s *Server) breakerMetrics(w http.ResponseWriter, r *http.Request) {
	m, err := s.svc.BreakerMetrics(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) ruleDistribution(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.RuleDistribution())
}

func (s *Server) topRejectedResources(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	httpx.OK(w, s.svc.TopRejectedResources(n))
}

func (s *Server) breakerEventDistribution(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.BreakerEventDistribution())
}

func (s *Server) topQuotaUsage(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	httpx.OK(w, s.svc.TopQuotaUsage(n))
}

func (s *Server) quotaReport(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.QuotaReport())
}

func (s *Server) adviseThresholds(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.AdviseThresholds())
}

func (s *Server) resourceStatuses(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ResourceStatuses())
}

type simulateRequest struct {
	Resource string `json:"resource"`
	Key      string `json:"key"`
	Count    int    `json:"count"`
}

func (s *Server) simulateTraffic(w http.ResponseWriter, r *http.Request) {
	var req simulateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.Count <= 0 {
		req.Count = 100
	}
	httpx.OK(w, s.svc.SimulateTraffic(req.Resource, req.Key, req.Count))
}

func (s *Server) breakerHealthSummary(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.BreakerHealthSummary())
}

func (s *Server) breakerRecoveries(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.BreakerRecoveries())
}

func (s *Server) ruleStatusReport(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.RuleStatusReport())
}

func (s *Server) resourceGroups(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ResourceGroups())
}

func (s *Server) decisionTrends(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.DecisionTrends())
}

func (s *Server) topRejectedKeys(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	httpx.OK(w, s.svc.TopRejectedKeys(n))
}

func (s *Server) reasonDistribution(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ReasonDistribution())
}

func (s *Server) alertAnalysis(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.AlertAnalysis())
}

func (s *Server) breakerTimeline(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.BreakerEventTimeline(r.PathValue("id")))
}

func (s *Server) breakerEventAnalysis(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.EventSequenceAnalysis(r.PathValue("id")))
}

func (s *Server) healthReport(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.HealthReport())
}

func (s *Server) quotaAnalysis(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.QuotaAnalysis())
}

func (s *Server) ruleRecommendations(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.RuleRecommendations())
}

func (s *Server) configReport(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ConfigReport())
}

func (s *Server) optimizationReport(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.OptimizationReport())
}

func (s *Server) ruleComplexityReport(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.RuleComplexityReport())
}

func (s *Server) breakerTuningSuggestions(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.BreakerTuningSuggestions())
}

func (s *Server) breakerEffectiveness(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.BreakerEffectiveness())
}

func (s *Server) trafficForecasts(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.TrafficForecasts())
}

func (s *Server) ruleStatsReport(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.RuleStatsReport())
}

func (s *Server) disabledRuleResources(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.DisabledRuleResources())
}

func (s *Server) resourceInventory(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ResourceInventory())
}

func (s *Server) systemSummary(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.SystemSummary())
}

func (s *Server) rulePresets(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, service.DefaultRulePresets())
}

func (s *Server) seedDefaultRules(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, map[string]int{"created": s.svc.SeedDefaultRules()})
}
