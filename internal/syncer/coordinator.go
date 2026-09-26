package syncer

import (
	"fmt"
	"sync"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

type Coordinator struct {
	mu      sync.Mutex
	state   *State
	windows map[model.EventID]*Window
}

func NewCoordinator(state *State) *Coordinator {
	return &Coordinator{state: state, windows: make(map[model.EventID]*Window)}
}

func (coordinator *Coordinator) Open(eventID model.EventID, required int) error {
	window, err := NewWindow(eventID, required)
	if err != nil {
		return err
	}
	if window.Required() != required {
		return fmt.Errorf("sync window requirement changed during creation")
	}
	coordinator.mu.Lock()
	defer coordinator.mu.Unlock()
	coordinator.windows[eventID] = window
	coordinator.state.Store(window.Status())
	return nil
}

func (coordinator *Coordinator) Observe(eventID model.EventID, frame model.Frame) (Status, error) {
	coordinator.mu.Lock()
	window := coordinator.windows[eventID]
	coordinator.mu.Unlock()
	if window == nil {
		return Status{}, fmt.Errorf("sync window for event %s is not open", eventID)
	}
	window.Observe(frame)
	status := window.Status()
	coordinator.state.Store(status)
	return status, nil
}

func (coordinator *Coordinator) Ready(eventID model.EventID) bool {
	status, ok := coordinator.state.Get(eventID)
	return ok && status.Stable
}

func (coordinator *Coordinator) Status(eventID model.EventID) (Status, bool) {
	return coordinator.state.Get(eventID)
}

func (coordinator *Coordinator) Close(eventID model.EventID) {
	coordinator.mu.Lock()
	delete(coordinator.windows, eventID)
	coordinator.mu.Unlock()
	status, _ := coordinator.state.Get(eventID)
	coordinator.state.Reset(eventID, status.Source, status.SourceGeneration, time.Now(), "sync window closed")
}
