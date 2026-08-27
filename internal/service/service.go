package service

import (
	"sync"
	"time"

	"flowcontrol/internal/config"
	"flowcontrol/internal/store"
	"flowcontrol/pkg/logger"
)

// Service 聚合业务逻辑，依赖 Store 与配置。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
	lim   *limiter
}

// New 构造服务实例。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{
		store: st,
		log:   log,
		cfg:   cfg,
		lim:   newLimiter(),
	}
}

// limiter 维护限流算法的运行时状态。
type limiter struct {
	mu      sync.Mutex
	windows map[string]*windowState
	sliding map[string][]time.Time
	buckets map[string]*bucketState
}

type windowState struct {
	start time.Time
	count int
}

type bucketState struct {
	tokens     float64
	lastRefill time.Time
}

func newLimiter() *limiter {
	return &limiter{
		windows: make(map[string]*windowState),
		sliding: make(map[string][]time.Time),
		buckets: make(map[string]*bucketState),
	}
}
