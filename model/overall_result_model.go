package model

import (
	"sync/atomic"
	"time"
)

type OverallResultAtomic struct {
	NumSuccess atomic.Uint64
	NumFailed  atomic.Uint64
	NumTotal   atomic.Uint64
}

func (m *OverallResultAtomic) GetOverallResultModel() *OverallResult {
	return &OverallResult{
		NumSuccess: m.NumSuccess.Load(),
		NumFailed:  m.NumFailed.Load(),
		NumTotal:   m.NumTotal.Load(),
	}
}

type OverallResult struct {
	NumSuccess uint64
	NumFailed  uint64
	NumTotal   uint64
	StartAt    time.Time
	EndAt      time.Time
}
