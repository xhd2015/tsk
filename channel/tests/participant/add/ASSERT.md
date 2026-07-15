## Expected

- AddParticipant succeeds (`added=true`).
- `participants.jsonl` has `alice` and `bob` sorted.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreOK(t, resp.StoreErr)
	if !resp.Added {
		t.Fatal("expected added=true")
	}
	assertParticipantHandlesSorted(t, activeChannelDir(req, "team-ch"), []string{"alice", "bob"})
}
```