package syncer

import (
	"fmt"
	"sync"

	"github.com/wyw14/solo-035/internal/model"
)

type Window struct {
	mu       sync.Mutex
	required int
	status   Status
	started  bool
}

func NewWindow(eventID model.EventID, required int) (*Window, error) {
	if required < 2 {
		return nil, fmt.Errorf("stable window requires at least two frames")
	}
	return &Window{required: required, status: Status{EventID: eventID}}, nil
}

func (window *Window) Observe(frame model.Frame) Status {
	window.mu.Lock()
	defer window.mu.Unlock()
	continuous := window.started &&
		window.status.Source == frame.Source &&
		window.status.SourceGeneration == frame.SourceGeneration &&
		window.status.LastSequence+1 == frame.Sequence &&
		window.timecodeContinuous(frame.Timecode)
	if !continuous {
		window.status.StableFrames = 0
		window.status.Stable = false
		window.status.Reason = "sequence restarted"
	}
	if frame.VideoLocked && frame.AudioLocked {
		window.status.StableFrames++
	} else {
		window.status.StableFrames = 0
		window.status.Reason = "media lock missing"
	}
	window.started = true
	window.status.Source = frame.Source
	window.status.SourceGeneration = frame.SourceGeneration
	window.status.LastSequence = frame.Sequence
	window.status.LastTimecode = frame.Timecode
	window.status.UpdatedAt = frame.CapturedAt.UTC()
	window.status.Stable = window.status.StableFrames >= window.required
	if window.status.Stable {
		window.status.Reason = "continuous media evidence"
	}
	return window.status
}

func (window *Window) timecodeContinuous(next int64) bool {
	if next == window.status.LastTimecode+1 {
		return true
	}
	if window.status.LastTimecode >= 2_591_999 && next == 0 {
		return false
	}
	return false
}

func (window *Window) Status() Status {
	window.mu.Lock()
	defer window.mu.Unlock()
	return window.status
}

func (window *Window) Required() int {
	window.mu.Lock()
	defer window.mu.Unlock()
	return window.required
}
