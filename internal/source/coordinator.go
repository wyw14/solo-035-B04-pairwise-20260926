package source

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

func (coordinator *Coordinator) EnsureSource(id model.SourceID) model.SignalStatus {
	if status, ok := coordinator.state.Status(id); ok {
		return status
	}
	return coordinator.state.Register(id, coordinator.clock())
}

func (coordinator *Coordinator) Switch(id model.SourceID) model.SignalStatus {
	coordinator.EnsureSource(id)
	return coordinator.state.BeginSwitch(id, coordinator.clock())
}

func (coordinator *Coordinator) Confirm(id model.SourceID, command model.CommandID, video bool, audio bool) error {
	if !coordinator.state.AcceptLock(id, command, video, audio, coordinator.clock()) {
		return fmt.Errorf("lock confirmation does not belong to current source command")
	}
	if !video || !audio {
		return fmt.Errorf("source %s is not fully locked", id)
	}
	return nil
}

func (coordinator *Coordinator) Ready(id model.SourceID) bool {
	status, ok := coordinator.state.Status(id)
	return ok && status.VideoLocked && status.AudioLocked
}

func (coordinator *Coordinator) Receive(frame model.Frame) error {
	coordinator.EnsureSource(frame.Source)
	if !coordinator.state.AddFrame(frame) {
		return fmt.Errorf("frame belongs to retired source generation")
	}
	return nil
}

func (coordinator *Coordinator) Retire(id model.SourceID) {
	coordinator.state.Retire(id)
}

func (coordinator *Coordinator) Snapshot() []model.SignalStatus {
	return coordinator.state.List()
}
