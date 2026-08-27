package store

import "flowcontrol/internal/model"

func (s *MemoryStore) CreateDecisionLog(d *model.DecisionLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.decisionLogs[d.ID] = d
	return nil
}

func (s *MemoryStore) ListDecisionLogs() []*model.DecisionLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DecisionLog, 0, len(s.decisionLogs))
	for _, d := range s.decisionLogs {
		list = append(list, d)
	}
	return list
}
