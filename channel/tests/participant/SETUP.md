# Scenario

**Feature**: channel Store participant membership

```
AddParticipant / RemoveParticipant with membership gate
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureStoreHelpersUsed()
	return nil
}
```