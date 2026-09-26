package source

import (
	"fmt"
	"sync"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

type Receiver struct {
	mu          sync.Mutex
	coordinator *Coordinator
	dropped     map[model.SourceID]int64
}

func NewReceiver(coordinator *Coordinator) *Receiver {
	return &Receiver{coordinator: coordinator, dropped: make(map[model.SourceID]int64)}
}

func (receiver *Receiver) Push(sourceID model.SourceID, generation uint64, sequence int64, timecode int64, capturedAt time.Time) (model.Frame, error) {
	frame := model.Frame{
		Source:           sourceID,
		SourceGeneration: generation,
		Sequence:         sequence,
		Timecode:         timecode,
		CapturedAt:       capturedAt.UTC(),
		VideoLocked:      true,
		AudioLocked:      true,
	}
	if err := receiver.coordinator.Receive(frame); err != nil {
		return model.Frame{}, err
	}
	return frame, nil
}

func (receiver *Receiver) RecordGap(sourceID model.SourceID, count int64) error {
	if count < 1 {
		return fmt.Errorf("gap must contain at least one missing frame")
	}
	receiver.mu.Lock()
	defer receiver.mu.Unlock()
	receiver.dropped[sourceID] += count
	return nil
}

func (receiver *Receiver) Dropped(sourceID model.SourceID) int64 {
	receiver.mu.Lock()
	defer receiver.mu.Unlock()
	return receiver.dropped[sourceID]
}
