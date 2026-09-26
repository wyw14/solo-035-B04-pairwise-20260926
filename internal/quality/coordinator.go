package quality

import (
	"fmt"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

type Coordinator struct {
	evidence *EvidenceSet
	clock    func() time.Time
}

func NewCoordinator(evidence *EvidenceSet, clock func() time.Time) *Coordinator {
	return &Coordinator{evidence: evidence, clock: clock}
}

func (coordinator *Coordinator) Record(eventID model.EventID, item model.AudioEvidence) error {
	if item.MeasuredAt.IsZero() {
		item.MeasuredAt = coordinator.clock().UTC()
	}
	return coordinator.evidence.Add(eventID, item)
}

func (coordinator *Coordinator) IndependentAudioSources(eventID model.EventID) int {
	unique := make(map[string]bool)
	for _, item := range coordinator.evidence.List(eventID) {
		origin := item.ProbeID
		if item.DerivedFrom != "" {
			origin = item.DerivedFrom
		}
		unique[origin] = true
	}
	return len(unique)
}

func (coordinator *Coordinator) Permit(eventID model.EventID, videoLocked bool, captionReady bool) (model.Readiness, error) {
	ready := model.Readiness{VideoLocked: videoLocked, CaptionReady: captionReady, CheckedAt: coordinator.clock().UTC()}
	ready.AudioReady = coordinator.IndependentAudioSources(eventID) >= 2
	ready.Ready = ready.VideoLocked && ready.AudioReady && ready.CaptionReady
	if !videoLocked {
		ready.Reason = "video is not locked"
	} else if !ready.AudioReady {
		ready.Reason = "independent audio evidence is missing"
	} else if !captionReady {
		ready.Reason = "caption chain is not ready"
	} else {
		ready.Reason = "all playout interlocks passed"
	}
	if !ready.Ready {
		return ready, fmt.Errorf("on-air permit denied: %s", ready.Reason)
	}
	return ready, nil
}

func (coordinator *Coordinator) Reset(eventID model.EventID) {
	coordinator.evidence.Clear(eventID)
}
