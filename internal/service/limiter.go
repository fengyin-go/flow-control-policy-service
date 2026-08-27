package service

import (
	"time"

	"flowcontrol/internal/model"
)

// allowByRule 根据规则算法判断请求是否放行。
func (s *Service) allowByRule(rule *model.RateLimitRule, key string) bool {
	s.lim.mu.Lock()
	defer s.lim.mu.Unlock()
	now := time.Now()
	switch rule.Algorithm {
	case model.AlgorithmSlidingWindow:
		return s.slidingAllow(rule, key, now)
	case model.AlgorithmTokenBucket:
		return s.bucketAllow(rule, key, now)
	default:
		return s.fixedWindowAllow(rule, key, now)
	}
}

// fixedWindowAllow 固定窗口限流。
func (s *Service) fixedWindowAllow(rule *model.RateLimitRule, key string, now time.Time) bool {
	id := rule.ID + ":" + key
	window := time.Duration(rule.WindowSec) * time.Second
	ws, ok := s.lim.windows[id]
	if !ok || now.Sub(ws.start) >= window {
		s.lim.windows[id] = &windowState{start: now, count: 1}
		return true
	}
	if ws.count >= rule.Limit {
		return false
	}
	ws.count++
	return true
}

// slidingAllow 滑动窗口限流（基于时间戳日志）。
func (s *Service) slidingAllow(rule *model.RateLimitRule, key string, now time.Time) bool {
	id := rule.ID + ":" + key
	window := time.Duration(rule.WindowSec) * time.Second
	times := s.lim.sliding[id]
	cutoff := now.Add(-window)
	kept := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	times = kept
	if len(times) >= rule.Limit {
		s.lim.sliding[id] = times
		return false
	}
	times = append(times, now)
	s.lim.sliding[id] = times
	return true
}

// bucketAllow 令牌桶限流。
func (s *Service) bucketAllow(rule *model.RateLimitRule, key string, now time.Time) bool {
	_ = key
	bs, ok := s.lim.buckets[rule.ID]
	capacity := float64(rule.Limit)
	rate := capacity / float64(rule.WindowSec)
	if !ok {
		bs = &bucketState{tokens: capacity, lastRefill: now}
		s.lim.buckets[rule.ID] = bs
	}
	elapsed := now.Sub(bs.lastRefill).Seconds()
	bs.tokens += elapsed * rate
	if bs.tokens > capacity {
		bs.tokens = capacity
	}
	bs.lastRefill = now
	if bs.tokens < 1 {
		return false
	}
	bs.tokens--
	return true
}
