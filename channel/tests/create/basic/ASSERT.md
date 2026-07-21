## Expected

- Store.Create succeeds.
- `channels/active/eng-alerts/channel.json` metadata only (`status: active`).
- `participants.jsonl` contains creator `alice` only.
- `channels/index/eng-alerts` is `active/eng-alerts`.
- `messages.jsonl` and `msg-counter` exist.

## Exit Code

- N/A (direct store call)

```go
import (
	"path/filepath"
	"os"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreOK(t, resp.StoreErr)

	id := "eng-alerts"
	dir := activeChannelDir(req, id)
	assertDirExists(t, dir)
	assertFileExists(t, filepath.Join(dir, "channel.json"))
	assertFileExists(t, filepath.Join(dir, "participants.jsonl"))
	assertFileExists(t, filepath.Join(dir, "messages.jsonl"))
	assertFileExists(t, filepath.Join(dir, "msg-counter"))

	ch := readChannelMetadata(t, dir)
	if ch.ID != id {
		t.Fatalf("id: got %q want %q", ch.ID, id)
	}
	if ch.Name != req.ChannelName {
		t.Fatalf("name: got %q want %q", ch.Name, req.ChannelName)
	}
	if ch.Status != "active" {
		t.Fatalf("status: got %q want active", ch.Status)
	}
	assertParticipantHandlesSorted(t, dir, []string{"alice"})

	idx := readChannelIndex(t, req, id)
	if idx != "active/"+id {
		t.Fatalf("index: got %q want active/%s", idx, id)
	}

	info, statErr := os.Stat(filepath.Join(dir, "messages.jsonl"))
	if statErr != nil {
		t.Fatalf("stat messages.jsonl: %v", statErr)
	}
	if info.Size() != 0 {
		t.Fatalf("messages.jsonl should be empty, got %d bytes", info.Size())
	}
}
```