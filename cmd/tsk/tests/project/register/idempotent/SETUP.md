# Scenario

**Feature**: re-register same project is idempotent (up-to-date or empty→fill)

```
register -> register again (same paths) -> already up to date
register (location empty on disk) -> fill location
```

## Intentional exclusions

- Conflicting non-empty path changes (covered by `register/basic` dup).
