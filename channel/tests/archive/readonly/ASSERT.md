## Expected

- Send on archived channel returns store error.
- No messages in archive dir.

## Errors

- Archived channel readonly.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreErr(t, resp.StoreErr)
	msgs := readMessagesJSONL(t, archiveChannelDir(req, "readonly-ch"))
	if len(msgs) != 0 {
		t.Fatalf("expected no messages, got %d", len(msgs))
	}
}
```