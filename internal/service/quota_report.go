package service

// QuotaReport 配额汇总报告。
type QuotaReport struct {
	TotalQuotas     int     `json:"total_quotas"`
	TotalAllowed    int     `json:"total_allowed"`
	TotalUsed       int     `json:"total_used"`
	ExhaustedQuotas int     `json:"exhausted_quotas"`
	OverallUsage    float64 `json:"overall_usage"`
}

// QuotaReport 返回配额汇总报告。
func (s *Service) QuotaReport() QuotaReport {
	report := QuotaReport{}
	for _, q := range s.store.ListQuotas() {
		report.TotalQuotas++
		report.TotalAllowed += q.Allowed
		report.TotalUsed += q.Used
		if q.Used >= q.Allowed {
			report.ExhaustedQuotas++
		}
	}
	if report.TotalAllowed > 0 {
		report.OverallUsage = float64(report.TotalUsed) / float64(report.TotalAllowed)
	}
	return report
}
