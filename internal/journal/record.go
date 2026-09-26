package journal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wyw14/solo-035/internal/model"
)

type Record struct {
	ID        string           `json:"id"`
	EventID   model.EventID    `json:"event_id"`
	Kind      string           `json:"kind"`
	Revision  model.RevisionID `json:"revision"`
	Fence     uint64           `json:"fence"`
	Payload   json.RawMessage  `json:"payload"`
	CreatedAt time.Time        `json:"created_at"`
}

func NewRecord(eventID model.EventID, kind string, revision model.RevisionID, fence uint64, payload any, now time.Time) (Record, error) {
	if !eventID.Valid() || kind == "" {
		return Record{}, fmt.Errorf("journal record identity is incomplete")
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return Record{}, fmt.Errorf("encode journal payload: %w", err)
	}
	return Record{ID: uuid.NewString(), EventID: eventID, Kind: kind, Revision: revision, Fence: fence, Payload: encoded, CreatedAt: now.UTC()}, nil
}

func (record Record) Decode(destination any) error {
	if len(record.Payload) == 0 {
		return fmt.Errorf("journal payload is empty")
	}
	if err := json.Unmarshal(record.Payload, destination); err != nil {
		return fmt.Errorf("decode journal payload: %w", err)
	}
	return nil
}
