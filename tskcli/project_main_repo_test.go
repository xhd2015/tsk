package tskcli

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGitMainRepoDirFromWorktree(t *testing.T) {
	main := t.TempDir()
	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run(main, "init")
	run(main, "config", "user.email", "t@t")
	run(main, "config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(main, "README"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(main, "add", "README")
	run(main, "commit", "-m", "init")

	wt := filepath.Join(t.TempDir(), "wt")
	run(main, "worktree", "add", wt, "HEAD")

	mainAbs, err := filepath.EvalSymlinks(main)
	if err != nil {
		mainAbs = main
	}
	wtAbs, err := filepath.EvalSymlinks(wt)
	if err != nil {
		wtAbs = wt
	}

	gotMain, err := gitMainRepoDir(mainAbs)
	if err != nil {
		t.Fatal(err)
	}
	gotWT, err := gitMainRepoDir(wtAbs)
	if err != nil {
		t.Fatal(err)
	}
	gotMain, _ = filepath.EvalSymlinks(gotMain)
	gotWT, _ = filepath.EvalSymlinks(gotWT)
	if gotMain != mainAbs {
		t.Fatalf("from main: got %q want %q", gotMain, mainAbs)
	}
	if gotWT != mainAbs {
		t.Fatalf("from worktree: got %q want main %q", gotWT, mainAbs)
	}
}
