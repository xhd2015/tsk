# Channel Store Tests

Store-level integration tests for `github.com/xhd2015/tsk/channel` and the
`github.com/xhd2015/tsk/channel/file` FileStore. Direct Go calls — no subprocess.

## Version
0.0.2

# DSN (Domain Specific Notion)

- **channel.Store** — provider-agnostic interface for channel lifecycle, messaging,
  and membership (`Create`, `List`, `Get`, `Archive`, `Delete`, `SendMessage`,
  `ListMessages`, `ListParticipants`, `AddParticipant`, `RemoveParticipant`).
- **file.FileStore** — file provider rooted at `TSK_HOME`; normalized on-disk layout
  matching the MySQL golden schema.
- **TSK_HOME** — storage root; tests isolate per leaf at `{WorkRoot}/.tsk`.
- **TSK_DATE** — optional env (`YYYY-MM-DD`) for deterministic timestamps; all tests
  set `TSK_DATE=2026-07-09`.
- **channel.json** — metadata only: `id`, `name`, `status`, `created_at`, `updated_at`
  (no embedded `participants`).
- **participants.jsonl** — one `{"handle","joined_at"}` per line, sorted by handle on
  write; creator is the only auto-joined participant on create (no `agent`).
- **messages.jsonl** — one message JSON per line; ids from `msg-counter` (flock).
- **index/<id>** — UTF-8 line `active/<id>` or `archive/<id>`.
- **tombstones/<id>.json** — `{"id","deleted_at"}` blocks id reuse after delete.
- **Membership gate** — non-participants cannot send or mutate membership; archived
  channels are readonly for mutations.
- **Creator identity** — tests default `Creator=alice` via `TSK_USER=alice`.

## Tree Overview

```
channel/tests
├── create/
│   ├── basic/                  # creator only in participants.jsonl; metadata-only channel.json
│   ├── duplicate/              # same id → error
│   ├── invalid-id/             # bad id format → error
│   └── tombstone-block/        # delete then recreate → error; tombstone remains
├── send/
│   └── basic/                  # messages.jsonl + msg-counter increment
├── list/
│   └── active-only/            # archived hidden by default
├── archive/
│   └── readonly/               # send blocked on archived channel
└── participant/
    ├── add/                    # add bob to roster
    ├── remove-self/            # creator removes another member; roster shrinks
    ├── last-participant/       # cannot remove sole member
    └── not-member/             # non-participant cannot add
```

## Test Index

| Leaf | Description |
|------|-------------|
| `create/basic` | Create → normalized layout; participants.jsonl has creator only |
| `create/duplicate` | Second create with same id → error |
| `create/invalid-id` | Invalid channel id → error |
| `create/tombstone-block` | Delete then recreate same id → error; tombstone persists |
| `send/basic` | Participant send → message in jsonl, counter bumped |
| `list/active-only` | List omits archived channels by default |
| `archive/readonly` | Send on archived channel → error |
| `participant/add` | Add bob → participants.jsonl includes bob sorted |
| `participant/remove-self` | Remove bob → alice remains |
| `participant/last-participant` | Remove sole alice → error |
| `participant/not-member` | Non-member add attempt → error |

## How to Run

```sh
doctest vet ./external/tsk-master-2026-07-14/channel/tests
doctest test ./external/tsk-master-2026-07-14/channel/tests

doctest test ./external/tsk-master-2026-07-14/channel/tests/create
doctest test ./external/tsk-master-2026-07-14/channel/tests/send
doctest test ./external/tsk-master-2026-07-14/channel/tests/list
doctest test ./external/tsk-master-2026-07-14/channel/tests/archive
doctest test ./external/tsk-master-2026-07-14/channel/tests/participant
```

```go
import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/tsk/channel"
	"github.com/xhd2015/tsk/channel/file"
)

type Request struct {
	WorkRoot    string
	TskHome     string
	Op          string
	Creator     string
	ChannelID   string
	ChannelName string
	MessageBody string
	Sender      string
	Handle      string
	ListAll     bool
}

type Response struct {
	Channel      *channel.Channel
	List         []channel.ListEntry
	Message      *channel.Message
	Messages     []channel.Message
	Participants []channel.Participant
	Added        bool
	StoreErr     error
}

func Run(t *testing.T, req *Request) (*Response, error) {
	restore := pushStoreEnv(req)
	defer restore()
	store := newFileStore(t, req)
	ctx := context.Background()
	resp := &Response{}

	switch req.Op {
	case "create":
		ch, err := store.Create(ctx, channel.CreateRequest{
			Name:    req.ChannelName,
			ID:      req.ChannelID,
			Creator: req.Creator,
		})
		resp.Channel = ch
		resp.StoreErr = err
	case "create_duplicate":
		_, err := store.Create(ctx, channel.CreateRequest{
			Name:    req.ChannelName,
			ID:      req.ChannelID,
			Creator: req.Creator,
		})
		if err != nil {
			resp.StoreErr = err
			return resp, nil
		}
		_, err = store.Create(ctx, channel.CreateRequest{
			Name:    req.ChannelName + " again",
			ID:      req.ChannelID,
			Creator: req.Creator,
		})
		resp.StoreErr = err
	case "create_invalid_id":
		_, err := store.Create(ctx, channel.CreateRequest{
			Name:    req.ChannelName,
			ID:      req.ChannelID,
			Creator: req.Creator,
		})
		resp.StoreErr = err
	case "create_tombstone_block":
		_, err := store.Create(ctx, channel.CreateRequest{
			Name:    req.ChannelName,
			ID:      req.ChannelID,
			Creator: req.Creator,
		})
		if err != nil {
			resp.StoreErr = err
			return resp, nil
		}
		if err := store.Delete(ctx, req.ChannelID); err != nil {
			resp.StoreErr = err
			return resp, nil
		}
		_, err = store.Create(ctx, channel.CreateRequest{
			Name:    req.ChannelName,
			ID:      req.ChannelID,
			Creator: req.Creator,
		})
		resp.StoreErr = err
	case "send":
		msg, err := store.SendMessage(ctx, channel.SendMessageRequest{
			ChannelID: req.ChannelID,
			Sender:    req.Sender,
			Body:      req.MessageBody,
		})
		resp.Message = msg
		resp.StoreErr = err
	case "list":
		entries, err := store.List(ctx, channel.ListOptions{All: req.ListAll})
		resp.List = entries
		resp.StoreErr = err
	case "archive_readonly":
		if err := store.Archive(ctx, req.ChannelID); err != nil {
			resp.StoreErr = err
			return resp, nil
		}
		_, err := store.SendMessage(ctx, channel.SendMessageRequest{
			ChannelID: req.ChannelID,
			Sender:    req.Sender,
			Body:      req.MessageBody,
		})
		resp.StoreErr = err
	case "participant_add":
		added, err := store.AddParticipant(ctx, channel.ParticipantChangeRequest{
			ChannelID: req.ChannelID,
			Handle:    req.Handle,
			Actor:     req.Sender,
		})
		resp.Added = added
		resp.StoreErr = err
	case "participant_remove":
		err := store.RemoveParticipant(ctx, channel.ParticipantChangeRequest{
			ChannelID: req.ChannelID,
			Handle:    req.Handle,
			Actor:     req.Sender,
		})
		resp.StoreErr = err
	case "participant_not_member":
		_, err := store.AddParticipant(ctx, channel.ParticipantChangeRequest{
			ChannelID: req.ChannelID,
			Handle:    req.Handle,
			Actor:     req.Sender,
		})
		resp.StoreErr = err
	default:
		return nil, fmt.Errorf("unknown op %q", req.Op)
	}
	return resp, nil
}

func pushStoreEnv(req *Request) func() {
	keys := []string{"TSK_HOME", "TSK_DATE"}
	saved := make(map[string]string, len(keys))
	for _, k := range keys {
		saved[k] = os.Getenv(k)
	}
	os.Setenv("TSK_HOME", req.TskHome)
	os.Setenv("TSK_DATE", "2026-07-09")
	return func() {
		for k, v := range saved {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}
}

func newFileStore(t *testing.T, req *Request) channel.Store {
	t.Helper()
	home := req.TskHome
	if home == "" {
		home = filepath.Join(req.WorkRoot, ".tsk")
	}
	st, err := file.NewStore(file.Options{Home: home})
	if err != nil {
		t.Fatalf("file.NewStore: %v", err)
	}
	return st
}
```