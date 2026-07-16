## Expected

- Non-member add returns store error naming actor and channel.
- `participants.jsonl` unchanged (`alice` only).

## Errors

- Membership gate blocks non-participant.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreErr(t, resp.StoreErr)
	want := `"carol" is not a participant in channel "private-ch"`
	if !strings.Contains(resp.StoreErr.Error(), want) {
		t.Fatalf("store err: got %q want substring %q", resp.StoreErr.Error(), want)
	}
	assertParticipantHandlesSorted(t, activeChannelDir(req, "private-ch"), []string{"alice"})
}
```
