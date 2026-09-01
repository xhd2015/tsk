package tskcli

import "testing"

func TestFormatTaskStageExtra(t *testing.T) {
	if got := formatTaskStageExtra("create", true); got != "" {
		t.Fatalf("color on create: got %q want empty", got)
	}
	if got := formatTaskStageExtra("done", true); got != "" {
		t.Fatalf("color on done: got %q want empty", got)
	}
	if got := formatTaskStageExtra("create", false); got != "  (create)" {
		t.Fatalf("color off create: got %q", got)
	}
	if got := formatTaskStageExtra("done", false); got != "  (done)" {
		t.Fatalf("color off done: got %q", got)
	}
}

func TestTaskStageStyle(t *testing.T) {
	if got := taskStageStyle("create", true); got != "" {
		t.Fatalf("create style: got %q", got)
	}
	if got := taskStageStyle("done", true); got != ansiGray+ansiStrikethrough {
		t.Fatalf("done style: got %q", got)
	}
	if got := taskStageStyle("in_process", true); got != ansiCyan {
		t.Fatalf("in_process style: got %q", got)
	}
	if got := taskStageStyle("done", false); got != "" {
		t.Fatalf("color off should not style: got %q", got)
	}
}

func TestFormatTaskLeafNameAlignment(t *testing.T) {
	// max width for [100] is 5
	got80 := formatTaskLeafName(80, "A", "a", 5)
	got100 := formatTaskLeafName(100, "C", "c", 5)
	got90 := formatTaskLeafName(90, "D", "d", 5)
	if got80 != "[80]  A" || got100 != "[100] C" || got90 != "[90]  D" {
		t.Fatalf("alignment:\n  %q\n  %q\n  %q", got80, got100, got90)
	}
}

func TestFormatProjectOriginExtra(t *testing.T) {
	if got := formatProjectOriginExtra("", true); got != "" {
		t.Fatalf("empty origin: %q", got)
	}
	if got := formatProjectOriginExtra("github.com/xhd2015/kck", false); got != "  github.com/xhd2015/kck" {
		t.Fatalf("plain: %q", got)
	}
	want := "  " + ansiGray + "github.com/xhd2015/kck" + ansiReset
	if got := formatProjectOriginExtra("github.com/xhd2015/kck", true); got != want {
		t.Fatalf("color: got %q want %q", got, want)
	}
}

func TestTruncateRunes(t *testing.T) {
	if got := truncateRunes("hello", 10); got != "hello" {
		t.Fatalf("short: %q", got)
	}
	runes := make([]rune, 600)
	for i := range runes {
		runes[i] = 'a'
	}
	got := truncateRunes(string(runes), 512)
	gr := []rune(got)
	if len(gr) != 512 {
		t.Fatalf("len=%d want 512", len(gr))
	}
	if gr[len(gr)-1] != '…' {
		t.Fatalf("expected ellipsis suffix")
	}
}
