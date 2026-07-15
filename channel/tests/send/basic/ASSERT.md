## Expected

- SendMessage succeeds; message id 1.
- `messages.jsonl` has one line; `msg-counter` is 1.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreOK(t, resp.StoreErr)
	if resp.Message == nil || resp.Message.ID != 1 {
		t.Fatalf("message id: got %+v want id 1", resp.Message)
	}
	dir := activeChannelDir(req, "send-ch")
	msgs := readMessagesJSONL(t, dir)
	if len(msgs) != 1 {
		t.Fatalf("messages: got %d want 1", len(msgs))
	}
	if msgs[0].Sender != "alice" || msgs[0].Body != "hello team" {
		t.Fatalf("message: got %+v", msgs[0])
	}
	if n := readMsgCounter(t, dir); n != 1 {
		t.Fatalf("msg-counter: got %d want 1", n)
	}
}
```