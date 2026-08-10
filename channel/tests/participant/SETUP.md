# Scenario

**Feature**: channel Store participant membership

```
AddParticipant / RemoveParticipant with membership gate
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureStoreHelpersUsed()
	return nil
}
```