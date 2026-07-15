## Expected

- RemoveParticipant succeeds.
- `participants.jsonl` has `alice` only.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreOK(t, resp.StoreErr)
	assertParticipantHandlesSorted(t, activeChannelDir(req, "team-ch"), []string{"alice"})
}
```