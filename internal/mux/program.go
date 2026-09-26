package mux

import (
	"fmt"
	"time"

	"github.com/wyw14/solo-035/internal/model"
)

type Program struct {
	EventID    model.EventID    `json:"event_id"`
	Revision   model.RevisionID `json:"revision"`
	Encoder    model.LeaseID    `json:"encoder_lease"`
	VideoPID   int              `json:"video_pid"`
	AudioPID   int              `json:"audio_pid"`
	CaptionPID int              `json:"caption_pid"`
	Bitrate    int              `json:"bitrate"`
	CreatedAt  time.Time        `json:"created_at"`
}

func NewProgram(eventID model.EventID, revision model.RevisionID, encoder model.LeaseID, bitrate int, now time.Time) (Program, error) {
	if !eventID.Valid() || !revision.Valid() || !encoder.Valid() {
		return Program{}, fmt.Errorf("program identities are incomplete")
	}
	if bitrate < 1000 {
		return Program{}, fmt.Errorf("program bitrate is below operational floor")
	}
	return Program{
		EventID:    eventID,
		Revision:   revision,
		Encoder:    encoder,
		VideoPID:   256,
		AudioPID:   257,
		CaptionPID: 258,
		Bitrate:    bitrate,
		CreatedAt:  now.UTC(),
	}, nil
}

func (program Program) Validate() error {
	if program.VideoPID == program.AudioPID || program.VideoPID == program.CaptionPID || program.AudioPID == program.CaptionPID {
		return fmt.Errorf("program elementary streams share a pid")
	}
	if program.Bitrate < 1000 {
		return fmt.Errorf("program bitrate is invalid")
	}
	return nil
}

func (program Program) WithBitrate(bitrate int) (Program, error) {
	program.Bitrate = bitrate
	if err := program.Validate(); err != nil {
		return Program{}, err
	}
	return program, nil
}
