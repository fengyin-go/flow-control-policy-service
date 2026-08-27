package service

import (
	"sort"
	"time"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/idgen"
)

// CreateCircuitBreaker 创建熔断器。
func (s *Service) CreateCircuitBreaker(input model.CircuitBreaker) (*model.CircuitBreaker, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	c := &model.CircuitBreaker{
		ID:               idgen.Hex(),
		Name:             input.Name,
		Resource:         input.Resource,
		FailureThreshold: input.FailureThreshold,
		SuccessThreshold: input.SuccessThreshold,
		TimeoutMs:        input.TimeoutMs,
		OpenDurationSec:  input.OpenDurationSec,
		State:            model.BreakerClosed,
		FailCount:        0,
		SuccessCount:     0,
		LastStateChange:  now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := s.store.CreateCircuitBreaker(c); err != nil {
		return nil, err
	}
	return c, nil
}

// GetCircuitBreaker 按 ID 获取熔断器。
func (s *Service) GetCircuitBreaker(id string) (*model.CircuitBreaker, error) {
	return s.store.GetCircuitBreaker(id)
}

// ListCircuitBreakers 返回全部熔断器。
func (s *Service) ListCircuitBreakers() []*model.CircuitBreaker {
	list := s.store.ListCircuitBreakers()
	sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.Before(list[j].CreatedAt) })
	return list
}

// UpdateCircuitBreaker 更新熔断器参数。
func (s *Service) UpdateCircuitBreaker(id string, input model.CircuitBreaker) (*model.CircuitBreaker, error) {
	existing, err := s.store.GetCircuitBreaker(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.FailureThreshold = input.FailureThreshold
	existing.SuccessThreshold = input.SuccessThreshold
	existing.TimeoutMs = input.TimeoutMs
	existing.OpenDurationSec = input.OpenDurationSec
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCircuitBreaker(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteCircuitBreaker 删除熔断器。
func (s *Service) DeleteCircuitBreaker(id string) error {
	return s.store.DeleteCircuitBreaker(id)
}

// CheckBreaker 判断熔断器是否放行（open 超时后自动转 half_open）。
func (s *Service) CheckBreaker(resource string) (bool, error) {
	b, err := s.store.GetCircuitBreakerByResource(resource)
	if err != nil {
		return true, nil
	}
	switch b.State {
	case model.BreakerClosed, model.BreakerHalfOpen:
		return true, nil
	case model.BreakerOpen:
		if time.Since(b.LastStateChange) >= time.Duration(b.OpenDurationSec)*time.Second {
			b.State = model.BreakerHalfOpen
			b.SuccessCount = 0
			b.LastStateChange = time.Now()
			b.UpdatedAt = time.Now()
			_ = s.store.UpdateCircuitBreaker(b)
			s.emitBreakerEvent(b.ID, model.BreakerEventHalfOpen, "进入半开状态")
			return true, nil
		}
		return false, nil
	default:
		return true, nil
	}
}

// RecordSuccess 记录一次成功调用。
func (s *Service) RecordSuccess(resource string) error {
	b, err := s.store.GetCircuitBreakerByResource(resource)
	if err != nil {
		return nil
	}
	switch b.State {
	case model.BreakerClosed:
		b.FailCount = 0
	case model.BreakerHalfOpen:
		b.SuccessCount++
		if b.SuccessCount >= b.SuccessThreshold {
			b.State = model.BreakerClosed
			b.FailCount = 0
			b.SuccessCount = 0
			b.LastStateChange = time.Now()
			s.emitBreakerEvent(b.ID, model.BreakerEventClosed, "半开状态连续成功，熔断器恢复关闭")
		}
	}
	b.UpdatedAt = time.Now()
	return s.store.UpdateCircuitBreaker(b)
}

// RecordFailure 记录一次失败调用。
func (s *Service) RecordFailure(resource string) error {
	b, err := s.store.GetCircuitBreakerByResource(resource)
	if err != nil {
		return nil
	}
	switch b.State {
	case model.BreakerClosed:
		b.FailCount++
		if b.FailCount >= b.FailureThreshold {
			b.State = model.BreakerOpen
			b.LastStateChange = time.Now()
			s.emitBreakerEvent(b.ID, model.BreakerEventOpened, "失败次数达到阈值，熔断打开")
		}
	case model.BreakerHalfOpen:
		b.State = model.BreakerOpen
		b.LastStateChange = time.Now()
		s.emitBreakerEvent(b.ID, model.BreakerEventOpened, "半开状态再次失败，熔断重新打开")
	}
	b.UpdatedAt = time.Now()
	return s.store.UpdateCircuitBreaker(b)
}
