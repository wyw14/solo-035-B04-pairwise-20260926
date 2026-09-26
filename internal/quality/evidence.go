package quality

import (
	"fmt"
	"sort"
	"sync"

	"github.com/wyw14/solo-035/internal/model"
)

type EvidenceSet struct {
	mu      sync.RWMutex
	byEvent map[model.EventID][]model.AudioEvidence
}

func NewEvidenceSet() *EvidenceSet {
	return &EvidenceSet{byEvent: make(map[model.EventID][]model.AudioEvidence)}
}

func (set *EvidenceSet) Add(eventID model.EventID, evidence model.AudioEvidence) error {
	if evidence.ProbeID == "" || evidence.Kind == "" {
		return fmt.Errorf("audio evidence lacks source identity")
	}
	set.mu.Lock()
	defer set.mu.Unlock()
	set.byEvent[eventID] = append(set.byEvent[eventID], evidence)
	return nil
}

func (set *EvidenceSet) List(eventID model.EventID) []model.AudioEvidence {
	set.mu.RLock()
	defer set.mu.RUnlock()
	items := append([]model.AudioEvidence(nil), set.byEvent[eventID]...)
	sort.Slice(items, func(i, j int) bool { return items[i].MeasuredAt.Before(items[j].MeasuredAt) })
	return items
}

func (set *EvidenceSet) Clear(eventID model.EventID) {
	set.mu.Lock()
	defer set.mu.Unlock()
	delete(set.byEvent, eventID)
}
