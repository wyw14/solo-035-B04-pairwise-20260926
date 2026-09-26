package quality

import (
	"fmt"

	"github.com/wyw14/solo-035/internal/model"
)

type Thresholds struct {
	MinimumLoudness  float64 `json:"minimum_loudness"`
	MaximumLoudness  float64 `json:"maximum_loudness"`
	MaximumAgeMillis int64   `json:"maximum_age_millis"`
}

type Evaluator struct {
	thresholds Thresholds
}

func NewEvaluator(thresholds Thresholds) (*Evaluator, error) {
	if thresholds.MinimumLoudness >= thresholds.MaximumLoudness {
		return nil, fmt.Errorf("quality thresholds are invalid")
	}
	if thresholds.MaximumAgeMillis < 1 {
		return nil, fmt.Errorf("quality evidence age must be positive")
	}
	return &Evaluator{thresholds: thresholds}, nil
}

func (evaluator *Evaluator) Accept(item model.AudioEvidence) error {
	if item.Value < evaluator.thresholds.MinimumLoudness || item.Value > evaluator.thresholds.MaximumLoudness {
		return fmt.Errorf("audio evidence is outside configured range")
	}
	if item.ProbeID == "" {
		return fmt.Errorf("audio evidence has no probe identity")
	}
	return nil
}

func (evaluator *Evaluator) Thresholds() Thresholds {
	return evaluator.thresholds
}
