# Scenario

**Feature**: cannot remove the last remaining participant

```
only agent+alice -> alice cannot remove self if that violates >=1 rule;
remove bob when only agent remains -> error
```

## Steps

1. Create; remove alice and bob until one left; attempt remove last.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Solo", "solo-ch")
	addParticipant(t, req, "solo-ch", "bob")
	runTskOK(t, req, "channel", "participant", "remove", "--channel-id", "solo-ch", "bob")
	runTskOK(t, req, "channel", "participant", "remove", "--channel-id", "solo-ch", "alice")
	req.Args = []string{"channel", "participant", "remove", "--channel-id", "solo-ch", "agent"}
	return nil
}
```
