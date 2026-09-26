package syncer

import (
	"fmt"

	"github.com/wyw14/solo-035/internal/model"
)

type Recovery struct {
	EventID          model.EventID  `json:"event_id"`
	Source           model.SourceID `json:"source"`
	SourceGeneration uint64         `json:"source_generation"`
	LastSequence     int64          `json:"last_sequence"`
	LastTimecode     int64          `json:"last_timecode"`
}

func (state *State) ExportRecovery() []Recovery {
	statuses := state.List()
	items := make([]Recovery, 0, len(statuses))
	for _, status := range statuses {
		items = append(items, Recovery{
			EventID:          status.EventID,
			Source:           status.Source,
			SourceGeneration: status.SourceGeneration,
			LastSequence:     status.LastSequence,
			LastTimecode:     status.LastTimecode,
		})
	}
	return items
}

func (state *State) RestoreRecovery(items []Recovery) error {
	for _, item := range items {
		if !item.EventID.Valid() || !item.Source.Valid() || item.SourceGeneration == 0 {
			return fmt.Errorf("invalid sync recovery record")
		}
		state.Store(Status{
			EventID:          item.EventID,
			Source:           item.Source,
			SourceGeneration: item.SourceGeneration,
			LastSequence:     item.LastSequence,
			LastTimecode:     item.LastTimecode,
			Reason:           "restored pending fresh evidence",
		})
	}
	return nil
}
