package mux

import (
	"fmt"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

type Coordinator struct {
	state *State
	clock func() time.Time
}

func NewCoordinator(state *State, clock func() time.Time) *Coordinator {
	return &Coordinator{state: state, clock: clock}
}

func (coordinator *Coordinator) Configure(event model.LiveEvent, lease model.EncoderLease, bitrate int) (Program, error) {
	if lease.Owner != event.ID {
		return Program{}, fmt.Errorf("encoder lease belongs to another event")
	}
	program, err := NewProgram(event.ID, event.Revision, lease.ID, bitrate, coordinator.clock())
	if err != nil {
		return Program{}, err
	}
	if err := coordinator.state.Put(program); err != nil {
		return Program{}, err
	}
	return program, nil
}

func (coordinator *Coordinator) Activate(eventID model.EventID, revision model.RevisionID) error {
	if !coordinator.state.Activate(eventID, revision) {
		return fmt.Errorf("program revision is not current")
	}
	return nil
}

func (coordinator *Coordinator) Stop(eventID model.EventID) {
	coordinator.state.Deactivate(eventID)
}

func (coordinator *Coordinator) Programs() []Program {
	return coordinator.state.List()
}
