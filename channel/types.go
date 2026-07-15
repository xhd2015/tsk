package channel

// Channel is channel metadata (no embedded participants).
type Channel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Participant is one roster entry.
type Participant struct {
	Handle   string `json:"handle"`
	JoinedAt string `json:"joined_at"`
}

// Message is one channel message.
type Message struct {
	ID        int    `json:"id"`
	Sender    string `json:"sender"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

// ListEntry is a channel summary for list output.
type ListEntry struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// CreateRequest creates a new channel.
type CreateRequest struct {
	Name    string
	ID      string
	Creator string
}

// ListOptions controls channel listing.
type ListOptions struct {
	All bool
}

// SendMessageRequest sends a message to a channel.
type SendMessageRequest struct {
	ChannelID string
	Sender    string
	Body      string
}

// ParticipantChangeRequest adds or removes a participant.
type ParticipantChangeRequest struct {
	ChannelID string
	Handle    string
	Actor     string
}