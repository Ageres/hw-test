package storage

import (
	"github.com/google/uuid"
)

type ProcEvent struct {
	ID string `json:"id" binding:"required"`
}

func (p *ProcEvent) ValidateProcEvent() error {
	if p == nil {
		return ErrProcEventIsNil
	}
	err := uuid.Validate(p.ID)
	if err != nil {
		return NewSError("failed to validate proc event id", err)
	}
	return nil
}
