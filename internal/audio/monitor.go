package audio

import (
	"fmt"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

type Monitor struct {
	minimum float64
	maximum float64
}

func NewMonitor(minimum float64, maximum float64) (*Monitor, error) {
	if minimum >= maximum {
		return nil, fmt.Errorf("audio monitor range is invalid")
	}
	return &Monitor{minimum: minimum, maximum: maximum}, nil
}

func (monitor *Monitor) Measure(probeID string, waveform float64, at time.Time) (model.AudioEvidence, model.AudioEvidence, error) {
	if probeID == "" {
		return model.AudioEvidence{}, model.AudioEvidence{}, fmt.Errorf("probe id is required")
	}
	if waveform < monitor.minimum || waveform > monitor.maximum {
		return model.AudioEvidence{}, model.AudioEvidence{}, fmt.Errorf("waveform is outside operational range")
	}
	wave := model.AudioEvidence{ProbeID: probeID, Kind: "waveform", Value: waveform, MeasuredAt: at.UTC()}
	loudness := model.AudioEvidence{ProbeID: probeID, Kind: "integrated-loudness", Value: waveform - 3.0, MeasuredAt: at.UTC(), DerivedFrom: probeID}
	return wave, loudness, nil
}

func (monitor *Monitor) WithinRange(value float64) bool {
	return value >= monitor.minimum && value <= monitor.maximum
}
