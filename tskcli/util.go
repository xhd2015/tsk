package tskcli

import (
	"fmt"
	"strconv"
	"strings"
)

// currentCtx is set during Run for event recording from subcommands.
var currentCtx *invocationContext

func initRunCtx(ctx *invocationContext) {
	currentCtx = ctx
}

func parseID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("tsk: invalid task id %q", s)
	}
	return id, nil
}

func joinArgs(args []string) string {
	return strings.Join(args, " ")
}