package service

import (
	"sort"
	"strings"
)

// ResourceGroup 资源分组。
type ResourceGroup struct {
	Prefix   string `json:"prefix"`
	Count    int    `json:"count"`
	Rejected int64  `json:"rejected"`
	Total    int64  `json:"total"`
}

// ResourceGroups 按资源前缀分组统计流量。
func (s *Service) ResourceGroups() []ResourceGroup {
	groups := make(map[string]*ResourceGroup)
	seen := make(map[string]bool)
	// 汇总所有出现过的资源
	for _, r := range s.store.ListRateLimitRules() {
		seen[r.Resource] = true
	}
	for _, b := range s.store.ListCircuitBreakers() {
		seen[b.Resource] = true
	}
	for _, t := range s.store.ListTrafficStats() {
		seen[t.Resource] = true
	}
	for resource := range seen {
		prefix := resourcePrefix(resource)
		g, ok := groups[prefix]
		if !ok {
			g = &ResourceGroup{Prefix: prefix}
			groups[prefix] = g
		}
		g.Count++
		if t, err := s.store.GetTrafficStatsByResource(resource); err == nil {
			g.Total += t.TotalRequests
			g.Rejected += t.RejectedRequests
		}
	}
	out := make([]ResourceGroup, 0, len(groups))
	for _, g := range groups {
		out = append(out, *g)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Prefix < out[j].Prefix })
	return out
}

// resourcePrefix 提取资源前缀（形如 /api/x 取 /api）。
func resourcePrefix(resource string) string {
	if !strings.HasPrefix(resource, "/") {
		return resource
	}
	parts := strings.Split(strings.Trim(resource, "/"), "/")
	if len(parts) <= 1 {
		return "/" + parts[0]
	}
	return "/" + parts[0]
}
