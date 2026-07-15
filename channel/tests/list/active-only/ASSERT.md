## Expected

- List returns only `active-one`; excludes `archived-one`.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertStoreOK(t, resp.StoreErr)
	ids := make([]string, len(resp.List))
	for i, e := range resp.List {
		ids[i] = e.ID
	}
	hasActive := false
	for _, id := range ids {
		if id == "archived-one" {
			t.Fatalf("list should not include archived-one: %v", ids)
		}
		if id == "active-one" {
			hasActive = true
		}
	}
	if !hasActive {
		t.Fatalf("list missing active-one: %v", ids)
	}
}
```