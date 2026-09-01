package tskcli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/localbin"
	lessflags "github.com/xhd2015/less-flags"
)

// installCatalog maps installable shim names to argv prefixes forwarded by the
// generated #!/bin/sh wrapper (exec <argv...> "$@").
var installCatalog = map[string][]string{
	"pmark": {"tsk", "project", "add"},
}

func installFail(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if strings.HasPrefix(msg, "Error:") {
		return err
	}
	return fmt.Errorf("Error: %s", msg)
}

func runInstall(_ string, args []string) error {
	// Install targets ~/.local/bin (user home), not TSK_HOME.
	setCommand(currentCtx, "install", args)

	var dryRun bool
	remaining, err := lessflags.
		Bool("--dry-run", &dryRun).
		Help("-h,--help", installHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return installFail(err)
	}
	if len(remaining) == 0 {
		return installFail(fmt.Errorf("tsk install: name required (try: %s)", installAvailableNames()))
	}
	if len(remaining) > 1 {
		return installFail(fmt.Errorf("tsk install: unexpected arguments %q", remaining[1:]))
	}
	name := remaining[0]
	argv, ok := installCatalog[name]
	if !ok {
		return installFail(fmt.Errorf("tsk install: unknown name %q (available: %s)", name, installAvailableNames()))
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return installFail(fmt.Errorf("tsk install: resolve home: %w", err))
	}
	binDir, err := localbin.DefaultDir(userHome)
	if err != nil {
		return installFail(err)
	}
	body := localbin.ForwarderBody(argv)
	path := filepath.Join(binDir, name)
	display := displayHomePath(userHome, path)

	if dryRun {
		printDryRunWrapper(path, display, body)
		_ = localbin.EnsureOnPATH(localbin.EnsureOpts{
			Home:    userHome,
			DestDir: binDir,
			DryRun:  true,
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		})
		return nil
	}

	if _, err := localbin.WriteScript(binDir, name, body); err != nil {
		return installFail(fmt.Errorf("tsk install: %w", err))
	}
	fmt.Printf("installed %s\n", display)
	_ = localbin.EnsureOnPATH(localbin.EnsureOpts{
		Home:    userHome,
		DestDir: binDir,
		Stderr:  os.Stderr,
	})
	return nil
}

func printDryRunWrapper(path, display, body string) {
	existing, err := os.ReadFile(path)
	switch {
	case err == nil && string(existing) == body:
		fmt.Printf("[dry-run] skip: %s (already identical)\n", display)
	case err == nil:
		fmt.Printf("[dry-run] would overwrite %s\n", display)
	case os.IsNotExist(err):
		fmt.Printf("[dry-run] would write %s\n", display)
	default:
		fmt.Fprintf(os.Stderr, "[dry-run] warning: probe %s: %v\n", display, err)
		fmt.Printf("[dry-run] would write %s\n", display)
	}
}

func installAvailableNames() string {
	names := make([]string, 0, len(installCatalog))
	for n := range installCatalog {
		names = append(names, n)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func displayHomePath(home, abs string) string {
	home = filepath.Clean(home)
	abs = filepath.Clean(abs)
	if home != "" && (abs == home || strings.HasPrefix(abs, home+string(os.PathSeparator))) {
		return "~" + abs[len(home):]
	}
	return abs
}

func installHelp() string {
	return `Usage: tsk install [--dry-run] <name>

Install a convenience CLI wrapper into ~/.local/bin and ensure that
directory is on PATH (bash/zsh).

Available:
  pmark   → tsk project add

Flags:
  --dry-run   print what would change; write nothing
  -h, --help  show this help
`
}
