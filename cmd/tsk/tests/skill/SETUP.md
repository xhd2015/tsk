# Scenario

**Feature**: tsk multi-topic skill surface via `tsk skill`

```
# skillcmd SingleSkill + embedded docs TreeFS
User -> tsk skill [--list|--show|--install|--help] -> stdout/stderr/exit
# nested path/TOPIC.md topics under docs/; root docs/SKILL.md
```

## Preconditions

- Skill content is embedded (`docs/` + `skillcmd.SingleSkill`); no TSK_HOME fixture required for skill commands beyond normal leaf isolation.

## Steps

1. Descendant setups set `req.Args` for the skill action under test.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
