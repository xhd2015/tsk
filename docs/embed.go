// Package docs embeds the tsk multi-topic skill tree (Shape 3).
//
// Layout:
//
//	docs/SKILL.md                 — root skill index
//	docs/<path>/TOPIC.md          — nested topics (slash paths)
//
// Wire into skillcmd.SingleSkill with RootContent=SkillMD and TreeFS.
package docs

import "embed"

// Name is the skill directory / skillcmd skill name.
const Name = "tsk"

// SkillMD is the root SKILL.md content (index only; no install plumbing flags).
//
//go:embed SKILL.md
var SkillMD string

// TreeFS holds nested topic directories (path/TOPIC.md). Root SKILL.md is not
// in TreeFS so ListTreeTopics only enumerates real topics.
//
//go:embed overview
//go:embed add
//go:embed tree
//go:embed topic
//go:embed workflow
//go:embed note
//go:embed channel
//go:embed project
//go:embed install
//go:embed working-on-task
var TreeFS embed.FS
