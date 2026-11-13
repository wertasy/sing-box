package docgen

import (
	"testing"
)

func TestHeadingIndexesSkipsFences(t *testing.T) {
	lines := []string{
		"# Title",
		"```json",
		"# not-a-heading",
		"```",
		"## Real",
	}
	indexes := headingIndexes(lines)
	if len(indexes) != 2 || lines[indexes[0]] != "# Title" || lines[indexes[1]] != "## Real" {
		t.Fatalf("heading indexes = %v", indexes)
	}
}

func TestTwinHeadingByOrdinalAcrossLevels(t *testing.T) {
	enLines := []string{
		"# TLS",
		"### Inbound",
		"```bash",
		"# not-a-heading",
		"```",
		"### Fields",
		"#### enabled",
	}
	zhLines := []string{
		"# TLS 中文",
		"## 入站",
		"```bash",
		"# 不是标题",
		"```",
		"## 字段",
		"### enabled",
	}
	resolved, err := twinHeadingByOrdinal(enLines, "### Fields", zhLines)
	if err != nil {
		t.Fatal(err)
	}
	if resolved != "## 字段" {
		t.Fatalf("resolved = %q, want %q", resolved, "## 字段")
	}
}

func TestCommentProseStripsDocBlocks(t *testing.T) {
	dir := t.TempDir()
	writeProbeFile(t, dir, "probe.go", `package option

type Options struct {
	// List of users.
	//
	// @Doc users.username:
	//    Username.
	//
	// @Doc users.password:
	//    Password.
	Users []string `+"`json:\"users,omitempty\"`"+`
}
`)
	index, err := LoadASTIndex(dir)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := index.Flatten("Options", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("entries = %v", entries)
	}
	prose := entries[0].Prose
	for _, line := range prose {
		if line == "" {
			continue
		}
		if line != "List of users." {
			t.Fatalf("prose leaked documentation blocks: %#v", prose)
		}
	}
}
