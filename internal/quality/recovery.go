package quality

import (
	"fmt"

	"github.com/wyw14/solo-035/internal/model"
)

type RecoveryRecord struct {
	EventID model.EventID         `json:"event_id"`
	Items   []model.AudioEvidence `json:"items"`
}

func (set *EvidenceSet) Export(eventIDs []model.EventID) []RecoveryRecord {
	records := make([]RecoveryRecord, 0, len(eventIDs))
	for _, eventID := range eventIDs {
		records = append(records, RecoveryRecord{EventID: eventID, Items: set.List(eventID)})
	}
	return records
}

func (set *EvidenceSet) Restore(records []RecoveryRecord) error {
	for _, record := range records {
		if !record.EventID.Valid() {
			return fmt.Errorf("quality recovery record has invalid event")
		}
		set.Clear(record.EventID)
		for _, item := range record.Items {
			if err := set.Add(record.EventID, item); err != nil {
				return err
			}
		}
	}
	return nil
}
