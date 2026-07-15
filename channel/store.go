package channel

import "context"

// Store is the provider-agnostic channel storage interface.
type Store interface {
	Create(ctx context.Context, req CreateRequest) (*Channel, error)
	List(ctx context.Context, opts ListOptions) ([]ListEntry, error)
	Get(ctx context.Context, id string) (*Channel, error)
	Archive(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	SendMessage(ctx context.Context, req SendMessageRequest) (*Message, error)
	ListMessages(ctx context.Context, channelID string) ([]Message, error)
	ListParticipants(ctx context.Context, channelID string) ([]Participant, error)
	AddParticipant(ctx context.Context, req ParticipantChangeRequest) (bool, error)
	RemoveParticipant(ctx context.Context, req ParticipantChangeRequest) error
}