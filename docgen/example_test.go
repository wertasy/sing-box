package docgen

import (
	"testing"
)

func TestRenderExampleEntry(t *testing.T) {
	lines := renderExampleEntry("server", `"127.0.0.1"`, "")
	if len(lines) != 1 || lines[0] != `  "server": "127.0.0.1",` {
		t.Errorf("lines = %v", lines)
	}
	lines = renderExampleEntry("detour", `""`, "or {}")
	if len(lines) != 1 || lines[0] != `  "detour": "", // or {}` {
		t.Errorf("lines = %v", lines)
	}
	multi := "[\n  \"a\"\n]"
	lines = renderExampleEntry("host_key", multi, "")
	// Continuation lines shift to the key nesting level (2 spaces).
	want := []string{`  "host_key": [`, `    "a"`, `  ],`}
	if len(lines) != len(want) {
		t.Fatalf("lines = %v, want %v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("lines[%d] = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestQuoteJSON(t *testing.T) {
	if got, want := quoteJSON("ssh"), `"ssh"`; got != want {
		t.Errorf("quoteJSON = %q, want %q", got, want)
	}
}
