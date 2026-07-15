package channel

import "errors"

var (
	ErrNotFound          = errors.New("channel not found")
	ErrAlreadyExists     = errors.New("channel already exists")
	ErrArchived          = errors.New("channel archived")
	ErrNotParticipant    = errors.New("not a participant")
	ErrLastParticipant   = errors.New("cannot remove last participant")
	ErrParticipantAbsent = errors.New("participant not in channel")
	ErrInvalidID         = errors.New("invalid channel id")
	ErrInvalidHandle     = errors.New("invalid handle")
)