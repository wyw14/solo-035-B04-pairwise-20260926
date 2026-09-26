package playout

import (
	"testing"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

func stoppedAt() time.Time {
	return time.Unix(200, 0).UTC()
}

func TestStopUnknownEventDoesNotCreateOutput(t *testing.T) {
	state := NewState()

	if stopped := state.Stop(model.NewEventID(), stoppedAt()); stopped {
		t.Fatal("stop on an event that was never staged reported success")
	}
	if items := state.List(); len(items) != 0 {
		t.Fatalf("List() = %d outputs, want none for unknown event", len(items))
	}
	if trace := state.Trace(model.NewEventID()); len(trace) != 0 {
		t.Fatalf("trace entries created for unknown event: %v", trace)
	}
}

func TestStopKnownEventStopsOutputAndRecordsTrace(t *testing.T) {
	state := NewState()
	eventID := model.NewEventID()
	output := Output{
		EventID:  eventID,
		Revision: model.NewRevisionID(),
		Fence:    7,
		Source:   model.NewSourceID(),
		Segment:  "current-program",
	}
	if err := state.Stage(output); err != nil {
		t.Fatal(err)
	}
	now := time.Unix(100, 0).UTC()
	if err := state.Publish(eventID, output.Revision, output.Fence, now); err != nil {
		t.Fatal(err)
	}

	stopTime := stoppedAt()
	if stopped := state.Stop(eventID, stopTime); !stopped {
		t.Fatal("stop on a staged and published output reported it was missing")
	}

	items := state.List()
	if len(items) != 1 {
		t.Fatalf("List() = %d outputs, want the existing one retained", len(items))
	}
	stopped := items[0]
	if stopped.OnAir {
		t.Fatal("stopped output is still on air")
	}
	if !stopped.UpdatedAt.Equal(stopTime) {
		t.Fatalf("UpdatedAt = %s, want %s", stopped.UpdatedAt, stopTime)
	}
	if stopped.EventID != eventID || stopped.Revision != output.Revision || stopped.Source != output.Source {
		t.Fatal("stopped output lost its identity")
	}
	trace := state.Trace(eventID)
	if len(trace) != 3 || trace[2] != "stopped" {
		t.Fatalf("trace = %v, want staged/published/stopped", trace)
	}
}
