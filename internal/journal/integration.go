package journal

import (
	"fmt"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

type Integration struct {
	store *Store
	state *State
	clock func() time.Time
}

func NewIntegration(store *Store, state *State, clock func() time.Time) *Integration {
	return &Integration{store: store, state: state, clock: clock}
}

func (integration *Integration) Persist(eventID model.EventID, kind string, revision model.RevisionID, fence uint64, payload any) (Record, error) {
	record, err := NewRecord(eventID, kind, revision, fence, payload, integration.clock())
	if err != nil {
		return Record{}, err
	}
	if err := integration.store.Append(record); err != nil {
		return Record{}, fmt.Errorf("persist %s proof: %w", kind, err)
	}
	integration.state.Apply(record)
	return record, nil
}

func (integration *Integration) Recover() ([]Record, error) {
	records, err := integration.store.ReadAll()
	if err != nil {
		return nil, err
	}
	integration.state.Reset()
	for _, record := range records {
		integration.state.Apply(record)
	}
	return records, nil
}

func (integration *Integration) Records() []Record {
	return integration.state.Records()
}
