## Expected

- Recreate after delete returns store error.
- `channels/tombstones/eng-alerts.json` exists.
- No active dir or index for `eng-alerts`.

## Errors

- Tombstoned id cannot be recreated.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreErr(t, resp.StoreErr)
	assertFileExists(t, channelAbs(req, "tombstones/eng-alerts.json"))
	assertFileNotExists(t, activeChannelDir(req, "eng-alerts"))
	assertFileNotExists(t, channelAbs(req, "index/eng-alerts"))
	ts := readTombstone(t, req, "eng-alerts")
	if ts.ID != "eng-alerts" {
		t.Fatalf("tombstone id: got %q want eng-alerts", ts.ID)
	}
}
```