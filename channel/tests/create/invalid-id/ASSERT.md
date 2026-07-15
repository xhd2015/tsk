## Expected

- Store.Create returns error.
- No `channels/active/bad-id` directory.

## Errors

- Invalid id format.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreErr(t, resp.StoreErr)
	assertFileNotExists(t, activeChannelDir(req, "bad-id"))
}
```