## Expected

- Non-member add returns store error.
- `participants.jsonl` unchanged (`alice` only).

## Errors

- Membership gate blocks non-participant.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreErr(t, resp.StoreErr)
	assertParticipantHandlesSorted(t, activeChannelDir(req, "private-ch"), []string{"alice"})
}
```