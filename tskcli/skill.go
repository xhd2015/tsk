package tskcli

import (
	"fmt"

	"github.com/xhd2015/skills/skillcmd"
	"github.com/xhd2015/tsk/docs"
)

func singleSkill() *skillcmd.SingleSkill {
	return &skillcmd.SingleSkill{
		Name:        docs.Name,
		RootContent: docs.SkillMD,
		TreeFS:      docs.TreeFS,
		Usage:       "tsk skill --install",
		Help:        skillHelp(),
	}
}

func runSkill(home string, args []string) error {
	setCommand(currentCtx, "skill", args)
	if err := singleSkill().Handle(args); err != nil {
		return fail(fmt.Errorf("tsk skill: %w", err))
	}
	return nil
}

func skillHelp() string {
	return `Usage: tsk skill --show [--header] [<topic-path>]
       tsk skill <topic-path> --show [--header]
       tsk skill --install [OPTIONS] [<dir>]
       tsk skill --list

Show the root SKILL.md index or a nested topic (path/TOPIC.md).
Install copies SKILL.md and nested TOPIC.md topics into agent skill directories.
List prints the skill name and every available topic path.
--help also lists available topics (see below).

Examples:
  tsk skill --show
  tsk skill --show create
  tsk skill create --show
  tsk skill --list
  tsk skill --install --dry-run
  tsk skill --install --help

Options:
  --show [--header] [path]   Print skill or topic content (header-only with --header)
  --install [OPTIONS] [dir]  Install skill files (see --install --help)
  --list                     Print skill name and all topic paths
  -h, --help                 Show this help and available topics
`
}
