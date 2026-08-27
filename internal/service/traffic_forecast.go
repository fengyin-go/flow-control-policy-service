package service

import (
	"sort"
)

// TrafficForecast 流量预测。
type TrafficForecast struct {
	Resource        string  `json:"resource"`
	CurrentTotal    int64   `json:"current_total"`
	CurrentRejected int64   `json:"current_rejected"`
	ForecastTotal   int64   `json:"forecast_total"`
	GrowthRate      float64 `json:"growth_rate"`
	Action          string  `json:"action"`
}

// TrafficForecasts 基于当前流量做简单预测（假设 20% 增长）。
func (s *Service) TrafficForecasts() []TrafficForecast {
	out := make([]TrafficForecast, 0)
	for _, t := range s.store.ListTrafficStats() {
		f := TrafficForecast{
			Resource:        t.Resource,
			CurrentTotal:    t.TotalRequests,
			CurrentRejected: t.RejectedRequests,
			GrowthRate:      0.2,
		}
		f.ForecastTotal = int64(float64(t.TotalRequests) * 1.2)
		action := "keep"
		if rule, err := s.store.GetRateLimitRuleByResource(t.Resource); err == nil {
			if f.ForecastTotal > int64(rule.Limit) {
				action = "scale_up"
			}
		} else {
			action = "protect"
		}
		f.Action = action
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CurrentTotal > out[j].CurrentTotal })
	return out
}
