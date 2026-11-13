package docgen

import (
	"strings"
	"testing"
)

func TestNormalizeLine(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "inline marks",
			input:    "==Required==",
			expected: "Required",
		},
		{
			name:     "bold marks",
			input:    "**Deprecated**",
			expected: "Deprecated",
		},
		{
			name:     "code spans",
			input:    "One of `tcp` `udp`.",
			expected: "One of tcp udp.",
		},
		{
			name:     "links",
			input:    "see [TLS](/configuration/shared/tls/) for details",
			expected: "see TLS (/configuration/shared/tls/) for details",
		},
		{
			name:     "whitespace",
			input:    "a   b\n\tc",
			expected: "a b c",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := normalizeLine(testCase.input); got != testCase.expected {
				t.Fatalf("normalizeLine(%q) = %q, want %q", testCase.input, got, testCase.expected)
			}
		})
	}
}

func TestExtractOldSection(t *testing.T) {
	content := `# Trojan

### Structure

` + "```json\n{}\n```" + `

### Fields

#### server

==Required==

The server address.

#### password

==Required==

The password.

### Dial Fields

See [Dial Fields](/configuration/shared/dial/).
`
	lines := strings.Split(content, "\n")
	section, err := extractOldSection(lines, "### Fields", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(section.keys) != 2 {
		t.Fatalf("expected 2 keys, got %v", section.keys)
	}
	if section.keys[0] != "server" || section.keys[1] != "password" {
		t.Fatalf("unexpected keys: %v", section.keys)
	}
	prose := normalizeBlock(section.prose["password"])
	if prose != "Required\nThe password." {
		t.Fatalf("unexpected prose: %q", prose)
	}
	if _, err := extractOldSection(lines, "### Missing", ""); err == nil {
		t.Fatal("expected error for missing heading")
	}
}

func TestExtractOldSectionTable(t *testing.T) {
	content := `### Fields

| Key | Type |
| --- | --- |
| ` + "`enabled`" + ` | Bool |

Enabled.
`
	lines := strings.Split(content, "\n")
	section, err := extractOldSection(lines, "### Fields", "table")
	if err != nil {
		t.Fatal(err)
	}
	if len(section.keys) != 1 || section.keys[0] != "enabled" {
		t.Fatalf("unexpected keys: %v", section.keys)
	}
}

func TestExtractOldSectionLevelInsensitive(t *testing.T) {
	// Old-layout pages document fields one level below a level-2 section;
	// the locator must find "### Fields" text at any level.
	content := `# MASQUE Client

## Structure

` + "```json\n{}\n```" + `

## Fields

### server

The server address.

### token.mode

One of the modes.

## Dial Fields

See [Dial Fields](/configuration/shared/dial/).
`
	lines := strings.Split(content, "\n")
	section, err := extractOldSection(lines, "### Fields", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(section.keys) != 2 || section.keys[0] != "server" || section.keys[1] != "token.mode" {
		t.Fatalf("unexpected keys: %v", section.keys)
	}
	if section.levels["server"] != 3 || section.levels["token.mode"] != 3 {
		t.Fatalf("unexpected levels: %v", section.levels)
	}
	prose := normalizeBlock(section.prose["server"])
	if prose != "The server address." {
		t.Fatalf("unexpected prose: %q", prose)
	}
	if _, err := extractOldSection(lines, "## Missing", ""); err == nil {
		t.Fatal("expected error for missing heading")
	}
}

func TestRenderFieldBlock(t *testing.T) {
	block := renderFieldBlock("password", []string{"==Required==", "", "The Trojan password."}, []string{"- Deprecated."}, 4)
	if !strings.Contains(block, "#### password\n\n==Required==\n\nThe Trojan password.\n\n- Deprecated.") {
		t.Fatalf("unexpected block:\n%s", block)
	}
	block = renderFieldBlock("password", []string{"Use it."}, nil, 3)
	if block != "### password\n\nUse it." {
		t.Fatalf("unexpected level-3 block: %q", block)
	}
}

func TestReorderDescription(t *testing.T) {
	prose := []string{"@Since 1.13.0", "", "@Quote:", "", "    Only supported on Linux.", "", "Do not reserve a port.", "", "This allows reuse."}
	expected := []string{"Do not reserve a port.", "", "This allows reuse.", "", "@Since 1.13.0", "", "@Quote:", "", "    Only supported on Linux."}
	if got := reorderDescription(prose); strings.Join(got, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("reorderDescription = %q, want %q", got, expected)
	}
	// All markers move after the description in a fixed order: required,
	// indicator, since, deprecated-in, admonitions, tab.
	mixed := []string{"@Tab \"Example\"", "", "    body", "", "@Quote:", "", "    quote body", "", "@Since 1.2.0", "", "@Deprecated in 1.4.0:", "", "    Use the new field.", "", "@Indicator: Server only", "", "@Required", "", "Description."}
	wantOrder := []string{"Description.", "", "@Required", "", "@Indicator: Server only", "", "@Since 1.2.0", "", "@Deprecated in 1.4.0:", "", "    Use the new field.", "", "@Quote:", "", "    quote body", "", "@Tab \"Example\"", "", "    body"}
	if got := reorderDescription(mixed); strings.Join(got, "\n") != strings.Join(wantOrder, "\n") {
		t.Fatalf("reorderDescription order = %q, want %q", got, wantOrder)
	}
	if got := reorderDescription([]string{"@Since 1.0.0"}); strings.Join(got, "\n") != "@Since 1.0.0" {
		t.Fatalf("marker-only prose must stay unchanged: %q", got)
	}
}

func TestEnumValuesEqual(t *testing.T) {
	if !enumValuesEqual([]string{"tcp", "udp"}, "udp, tcp") {
		t.Fatal("expected equal sets")
	}
	if enumValuesEqual([]string{"tcp", "udp"}, "tcp") {
		t.Fatal("expected unequal sets")
	}
}

func TestRenderTagLineDefault(t *testing.T) {
	entries := []FlattenEntry{{Tags: FieldTags{Default: "22"}}}
	en, err := renderTagLine(entries, nil, nil, "en")
	if err != nil {
		t.Fatal(err)
	}
	if len(en) != 1 || en[0] != "`22` will be used by default." {
		t.Fatalf("en default = %v", en)
	}
	zh, err := renderTagLine(entries, nil, nil, "zh")
	if err != nil {
		t.Fatal(err)
	}
	if len(zh) != 1 || zh[0] != "默认使用 `22`。" {
		t.Fatalf("zh default = %v", zh)
	}
}

func TestExpandProse(t *testing.T) {
	testCases := []struct {
		name     string
		prose    []string
		lang     string
		expected []string
	}{
		{"required en", []string{"@Required"}, "en", []string{"==Required=="}},
		{"required zh", []string{"@Required"}, "zh", []string{"==必填=="}},
		{"indicator", []string{"@Indicator: Server only"}, "en", []string{"==Server only=="}},
		{"admonition titled", []string{`@Admonition question "Since sing-box 1.11.0"`}, "en", []string{`!!! question "Since sing-box 1.11.0"`}},
		{"admonition empty title", []string{`@Admonition note ""`}, "en", []string{`!!! note ""`}},
		{"admonition bare", []string{"@Admonition note"}, "en", []string{"!!! note"}},
		{"tab", []string{`@Tab "With auto_redirect"`}, "en", []string{`=== "With auto_redirect"`}},
		{"nested indented marker", []string{`    @Admonition quote ""`, "        body"}, "en", []string{`    !!! quote ""`, "        body"}},
		{"shorthand bare", []string{"@Note"}, "en", []string{"!!! note"}},
		{"shorthand titled", []string{"@Note: Some title"}, "en", []string{`!!! note "Some title"`}},
		{"shorthand colon titleless", []string{"@Quote:"}, "en", []string{`!!! quote ""`}},
		{"shorthand empty title", []string{`@Quote: ""`}, "en", []string{`!!! quote ""`}},
		{"shorthand since en", []string{"@Since 1.11.0"}, "en", []string{`!!! question "Since sing-box 1.11.0"`}},
		{"shorthand since zh", []string{"@Since 1.11.0"}, "zh", []string{`!!! question "自 sing-box 1.11.0 起"`}},
		{"deprecated in en", []string{"@Deprecated in 1.14.0:", "", "    Use the new field."}, "en", []string{`!!! failure "Deprecated in sing-box 1.14.0"`, "", "    Use the new field."}},
		{"deprecated in zh", []string{"@Deprecated in 1.14.0:", "", "    请使用新字段。"}, "zh", []string{`!!! failure "已在 sing-box 1.14.0 废弃"`, "", "    请使用新字段。"}},
		{"deprecated in indented", []string{"    @Deprecated in 1.2.0:", "", "        body"}, "en", []string{`    !!! failure "Deprecated in sing-box 1.2.0"`, "", "        body"}},
		{"shorthand indented", []string{"    @Warning"}, "en", []string{"    !!! warning"}},
		{"unknown at word untouched", []string{"@unknown"}, "en", []string{"@unknown"}},
		{"plain prose untouched", []string{"The server address."}, "en", []string{"The server address."}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			expanded := expandProse(testCase.prose, testCase.lang, testCase.lang)
			if strings.Join(expanded, "\n") != strings.Join(testCase.expected, "\n") {
				t.Fatalf("expandProse = %q, want %q", expanded, testCase.expected)
			}
		})
	}
}

func TestDefaultValueRetained(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		newProse string
		retained bool
	}{
		{"value kept", "`https://example.com/x.zip` will be used by default.", "Tags\n\n`https://example.com/x.zip`", true},
		{"value dropped", "`https://example.com/x.zip` will be used by default.", "No marker here.", false},
		{"no backticks", "This is used by default.", "This is used by default.", true},
		{"one of several kept", "Defaults to `80` or `443`.", "Defaults to `443` when empty.", true},
		{"none kept", "Defaults to `80` or `443`.", "Defaults to `8080` when empty.", false},
		{"generated default sentence", "Server port. 22 will be used if empty.", "`22` will be used by default.", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultValueRetained(tt.line, tt.newProse, tt.newProse); got != tt.retained {
				t.Fatalf("defaultValueRetained(%q, %q) = %v, want %v", tt.line, tt.newProse, got, tt.retained)
			}
		})
	}
}
