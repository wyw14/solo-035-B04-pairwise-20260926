package playout

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

type Integration struct {
	state   *State
	mu      sync.Mutex
	cancels map[model.EventID]context.CancelFunc
}

func NewIntegration(state *State) *Integration {
	return &Integration{state: state, cancels: make(map[model.EventID]context.CancelFunc)}
}

func (integration *Integration) StartWorker(parent context.Context, event model.LiveEvent, delay time.Duration, action func(context.Context) error) <-chan error {
	integration.mu.Lock()
	if cancel := integration.cancels[event.ID]; cancel != nil {
		cancel()
	}
	workerContext, cancel := context.WithCancel(parent)
	integration.cancels[event.ID] = cancel
	integration.mu.Unlock()
	result := make(chan error, 1)
	go func() {
		defer close(result)
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-workerContext.Done():
			result <- workerContext.Err()
		case <-timer.C:
			result <- action(workerContext)
		}
	}()
	return result
}

func (integration *Integration) Cancel(eventID model.EventID) error {
	integration.mu.Lock()
	defer integration.mu.Unlock()
	cancel := integration.cancels[eventID]
	if cancel == nil {
		return fmt.Errorf("playout worker for event %s does not exist", eventID)
	}
	cancel()
	delete(integration.cancels, eventID)
	return nil
}

func (integration *Integration) Active(eventID model.EventID) bool {
	integration.mu.Lock()
	defer integration.mu.Unlock()
	return integration.cancels[eventID] != nil
}
