# Scenario

**Feature**: list hides archived channels by default

```
create active-one + archived-one -> List (no All) -> active-one only
```

## Steps

1. Seed two channels; archive one.
2. List without All flag.

```go
func Setup(t *testing.T, req *Request) error {
	seedChannel(t, req, "Active One", "active-one")
	seedChannel(t, req, "Archived One", "archived-one")
	store := newFileStore(t, req)
	if err := store.Archive(context.Background(), "archived-one"); err != nil {
		t.Fatalf("archive: %v", err)
	}
	req.Op = "list"
	req.ListAll = false
	return nil
}
```