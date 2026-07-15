## Expected

- Remove last participant returns store error.
- `participants.jsonl` still has `alice`.

## Errors

- Cannot remove sole member.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreErr(t, resp.StoreErr)
	assertParticipantHandlesSorted(t, activeChannelDir(req, "solo-ch"), []string{"alice"})
}
```