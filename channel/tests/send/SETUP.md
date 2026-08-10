# Scenario

**Feature**: channel Store SendMessage

```
Participant -> Store.SendMessage -> messages.jsonl + msg-counter
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureStoreHelpersUsed()
	return nil
}
```