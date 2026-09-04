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

func runSkill(invk *invocation, args []string) error {
	invk.setCommand("skill", args)
	if err := singleSkill().Handle(args); err != nil {
		return fail(fmt.Errorf("tsk skill: %w", err))
	}
	return nil
}

func skillHelp() string {
	return `Usage: tsk skill --show [--header] [<topic-path>]
       tsk skill <topic-path> --show [--header]
       tsk skill --install [OPTIONS] [<dir>]
       tsk skill <topic-path> --install [OPTIONS] [<dir>]
       tsk skill --list

Show the root SKILL.md index or a nested topic (path/TOPIC.md).
Install copies SKILL.md and nested TOPIC.md topics into agent skill directories.
A leading <topic-path> before --install is not a destination (whole skill is
installed); pass --dir or <dir> for the target.
List prints the skill name and every available topic path.
--help also lists available topics (see below).

Examples:
  tsk skill --show
  tsk skill --show add
  tsk skill add --show
  tsk skill --show working-on-task
  tsk skill --list
  tsk skill --install --dry-run
  tsk skill --install --dir ~/skills --dry-run
  tsk skill working-on-task --install --dir ~/skills --dry-run
  tsk skill --install --help

Options:
  --show [--header] [path]   Print skill or topic content (header-only with --header)
  --install [OPTIONS] [dir]  Install skill files (see --install --help)
  --list                     Print skill name and all topic paths
  -h, --help                 Show this help and available topics
`
}
