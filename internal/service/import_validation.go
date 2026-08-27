package service

// SnapshotValidation 快照校验结果。
type SnapshotValidation struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

// ValidateSnapshot 校验快照引用完整性。
func (s *Service) ValidateSnapshot(snap Snapshot) SnapshotValidation {
	v := SnapshotValidation{Valid: true, Errors: []string{}}
	ruleIDs := make(map[string]bool)
	for _, r := range snap.Rules {
		ruleIDs[r.ID] = true
	}
	for _, t := range snap.TokenBuckets {
		if !ruleIDs[t.RuleID] {
			v.Errors = append(v.Errors, "令牌桶引用了不存在的规则 "+t.RuleID)
		}
	}
	for _, q := range snap.Quotas {
		if !ruleIDs[q.RuleID] {
			v.Errors = append(v.Errors, "配额引用了不存在的规则 "+q.RuleID)
		}
	}
	if len(v.Errors) > 0 {
		v.Valid = false
	}
	return v
}
