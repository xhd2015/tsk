---
name: tsk/channel
description: >-
  Slack-like channel spaces: create, send, messages, participants, archive.
---

# channel

Channels live under `TSK_HOME/channels/` (separate from the task tree).

```text
tsk channel create [--channel-id ID] [--user HANDLE] <name>
tsk channel list [--all] [--json]
tsk channel send --channel-id ID [--user HANDLE] <body…>
tsk channel messages --channel-id ID
tsk channel participant add|remove …
tsk channel archive|delete …
```

## Identity

Handle precedence: `--user` > `TSK_USER` > `$USER`. Format:
`^[a-z0-9][a-z0-9_-]{0,63}$`.

`--channel-id` and `--user` may appear after `channel` before the action
(parent peel); conflicting parent/leaf values error.

## Membership

Non-participants cannot send or mutate membership. Archived channels are
read-only for mutations; list/messages still work with `--all` where applicable.

## Agent guidance

Prefer task notes/progress for work tracking. Use channels when the user wants
a conversational space keyed by channel id, not a task id. For full flags, run
`tsk channel --help` and `tsk channel <action> --help`.
