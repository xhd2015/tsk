package tskcli

import "testing"

func TestNormalizeOriginURL(t *testing.T) {
	cases := []struct {
		in       string
		wantID   string
		wantName string
	}{
		{
			in:       "https://github.com/xhd2015/dot-pkgs",
			wantID:   "github.com/xhd2015/dot-pkgs",
			wantName: "dot-pkgs",
		},
		{
			in:       "https://github.com/xhd2015/my-private-working-note.git",
			wantID:   "github.com/xhd2015/my-private-working-note",
			wantName: "my-private-working-note",
		},
		{
			in:       "gitlab@git.example.com:acme/loan-service/credit_backend/code-lens/tools/widget-cli.git",
			wantID:   "git.example.com/acme/loan-service/credit_backend/code-lens/tools/widget-cli",
			wantName: "widget-cli",
		},
		{
			in:       "git@github.com:xhd2015/wrk.git",
			wantID:   "github.com/xhd2015/wrk",
			wantName: "wrk",
		},
		{
			in:       "https://git.example.com/acme/loan-service/credit_backend/widget_disburse_service.git",
			wantID:   "git.example.com/acme/loan-service/credit_backend/widget_disburse_service",
			wantName: "widget_disburse_service",
		},
		{
			in:       "ssh://git@github.com/xhd2015/agent-pro.git",
			wantID:   "github.com/xhd2015/agent-pro",
			wantName: "agent-pro",
		},
	}
	for _, tc := range cases {
		id, name, err := NormalizeOriginURL(tc.in)
		if err != nil {
			t.Fatalf("%q: %v", tc.in, err)
		}
		if id != tc.wantID || name != tc.wantName {
			t.Fatalf("%q: got (%q, %q) want (%q, %q)", tc.in, id, name, tc.wantID, tc.wantName)
		}
	}
}

func TestNormalizeOriginURLErrors(t *testing.T) {
	for _, in := range []string{"", "   ", "not-a-url"} {
		if _, _, err := NormalizeOriginURL(in); err == nil {
			t.Fatalf("%q: expected error", in)
		}
	}
}
