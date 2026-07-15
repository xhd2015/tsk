## Expected

- Second create returns store error.
- Original `channels/active/eng-alerts/` unchanged.

## Errors

- Duplicate channel id rejected.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreErr(t, resp.StoreErr)
	assertDirExists(t, activeChannelDir(req, "eng-alerts"))
}
```