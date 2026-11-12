package docgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderOriginBody(t *testing.T) {
	info := originSections["dial"]
	if got, want := renderOriginBody(info, "en"), "See [Dial Fields](/configuration/shared/dial/) for details."; got != want {
		t.Errorf("en body = %q, want %q", got, want)
	}
	if got, want := renderOriginBody(info, "zh"), "参阅 [拨号字段](/zh/configuration/shared/dial/) 了解详情。"; got != want {
		t.Errorf("zh body = %q, want %q", got, want)
	}
	listen := originSections["listen"]
	if got, want := renderOriginBody(listen, "en"), "See [Listen Fields](/configuration/shared/listen/) for details."; got != want {
		t.Errorf("listen en body = %q, want %q", got, want)
	}
}

func TestAuditIndexTables(t *testing.T) {
	variants := map[string]*SchemaNode{
		"alpha": {Properties: &SchemaNodeMap{}},
		"beta":  {Properties: &SchemaNodeMap{}},
	}
	table := indexTable{file: "outbound/index.md", registry: "Test", legacy: map[string]bool{"legacy-row": true}}
	schemaIndex := &SchemaIndex{Variants: map[string]map[string]*SchemaNode{"Test": variants}}

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "outbound"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTable := func(t *testing.T, rows ...string) {
		t.Helper()
		content := "| Type | Format |\n|------|--------|\n" + strings.Join(rows, "\n")
		if err := os.WriteFile(filepath.Join(dir, "outbound/index.md"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("complete", func(t *testing.T) {
		writeTable(t, "| `alpha` | [A](./a/) |", "| `beta` | [B](./b/) |", "| `legacy-row` | [L](./l/) |")
		if issues := auditTableForTest(schemaIndex, dir, table); len(issues) != 0 {
			t.Errorf("unexpected issues: %v", issues)
		}
	})

	t.Run("missing", func(t *testing.T) {
		writeTable(t, "| `alpha` | [A](./a/) |")
		issues := auditTableForTest(schemaIndex, dir, table)
		if len(issues) != 1 || !strings.Contains(issues[0].Error(), "index table missing type beta") {
			t.Errorf("issues = %v, want missing beta", issues)
		}
	})

	t.Run("unknown", func(t *testing.T) {
		writeTable(t, "| `alpha` | [A](./a/) |", "| `beta` | [B](./b/) |", "| `gamma` | [G](./g/) |")
		issues := auditTableForTest(schemaIndex, dir, table)
		if len(issues) != 1 || !strings.Contains(issues[0].Error(), "index table has unknown type gamma") {
			t.Errorf("issues = %v, want unknown gamma", issues)
		}
	})

	t.Run("legacy_allowed", func(t *testing.T) {
		writeTable(t, "| `alpha` | [A](./a/) |", "| `beta` | [B](./b/) |", "| `legacy-row` | [L](./l/) |", "| `other-legacy` | [O](./o/) |")
		issues := auditTableForTest(schemaIndex, dir, table)
		if len(issues) != 1 || !strings.Contains(issues[0].Error(), "unknown type other-legacy") {
			t.Errorf("issues = %v, want only other-legacy unknown", issues)
		}
	})
}

// auditTableForTest runs auditIndexTables with a single table override.
func auditTableForTest(schemaIndex *SchemaIndex, docsDir string, table indexTable) []error {
	saved := indexTables
	indexTables = []indexTable{table}
	defer func() { indexTables = saved }()
	return auditIndexTables(schemaIndex, docsDir)
}
