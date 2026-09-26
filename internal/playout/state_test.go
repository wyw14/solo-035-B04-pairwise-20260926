package playout

import (
	"testing"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

func TestStopUnknownEventDoesNotCreateOutput(t *testing.T) {
	state := NewState()
	state.Stop(model.NewEventID(), time.Now())
	if outputs := state.List(); len(outputs) != 0 {
		t.Fatalf("expected no outputs after stopping an unknown event, got %d", len(outputs))
	}
	if trace := state.Trace(model.NewEventID()); len(trace) != 0 {
		t.Fatalf("expected no trace for an unknown event, got %v", trace)
	}
}

func TestStopStagedOutputKeepsRecord(t *testing.T) {
	state := NewState()
	eventID := model.NewEventID()
	staged := Output{
		EventID:  eventID,
		Revision: model.NewRevisionID(),
		Fence:    7,
		Source:   model.NewSourceID(),
		Segment:  "segment-a",
	}
	if err := state.Stage(staged); err != nil {
		t.Fatalf("stage: %v", err)
	}
	if err := state.Publish(eventID, staged.Revision, staged.Fence, time.Now()); err != nil {
		t.Fatalf("publish: %v", err)
	}
	state.Stop(eventID, time.Now())
	outputs := state.List()
	if len(outputs) != 1 {
		t.Fatalf("expected the stopped output to remain listed, got %d outputs", len(outputs))
	}
	if outputs[0].EventID != eventID || !outputs[0].Revision.Valid() || !outputs[0].Source.Valid() {
		t.Fatalf("stopped output lost its identity: %+v", outputs[0])
	}
	if outputs[0].OnAir {
		t.Fatal("expected output to be off air after stop")
	}
	trace := state.Trace(eventID)
	if len(trace) == 0 || trace[len(trace)-1] != "stopped" {
		t.Fatalf("expected stop to be recorded in trace, got %v", trace)
	}
}
