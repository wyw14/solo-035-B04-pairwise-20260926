package backup

import (
	"fmt"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

type Decision struct {
	EventID     model.EventID `json:"event_id"`
	Reason      string        `json:"reason"`
	Severity    int           `json:"severity"`
	RequestedAt time.Time     `json:"requested_at"`
	Approved    bool          `json:"approved"`
}

func Evaluate(event model.LiveEvent, videoDrops int, audioDrops int, now time.Time) (Decision, error) {
	decision := Decision{EventID: event.ID, RequestedAt: now.UTC()}
	if event.Phase != model.PhaseOnAir && event.Phase != model.PhaseHolding {
		return decision, fmt.Errorf("event phase %s cannot fail over", event.Phase)
	}
	decision.Severity = videoDrops*3 + audioDrops*2
	if videoDrops > 0 {
		decision.Reason = "video continuity lost"
	} else if audioDrops > 1 {
		decision.Reason = "audio continuity lost"
	} else {
		decision.Reason = "signal remains within protection window"
	}
	decision.Approved = decision.Severity >= 3
	return decision, nil
}

func (decision Decision) Validate(event model.LiveEvent) error {
	if decision.EventID != event.ID {
		return fmt.Errorf("failover decision belongs to another event")
	}
	if decision.Approved && decision.Severity < 3 {
		return fmt.Errorf("failover severity is below approval threshold")
	}
	return nil
}
