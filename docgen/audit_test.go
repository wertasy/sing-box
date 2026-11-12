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

func TestRenderFieldBlock(t *testing.T) {
	block := renderFieldBlock("password", []string{"==Required==", "", "The Trojan password."}, []string{"- Deprecated."})
	if !strings.Contains(block, "#### password\n\n==Required==\n\nThe Trojan password.\n\n- Deprecated.") {
		t.Fatalf("unexpected block:\n%s", block)
	}
}

func TestReorderDescription(t *testing.T) {
	prose := []string{"@Since 1.13.0", "", "@Quote:", "", "    Only supported on Linux.", "", "Do not reserve a port.", "", "This allows reuse."}
	expected := []string{"Do not reserve a port.", "", "This allows reuse.", "", "@Since 1.13.0", "", "@Quote:", "", "    Only supported on Linux."}
	if got := reorderDescription(prose); strings.Join(got, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("reorderDescription = %q, want %q", got, expected)
	}
	// All markers move after the description in a fixed order: required,
	// indicator, since, admonitions, tab.
	mixed := []string{"@Tab \"Example\"", "", "    body", "", "@Quote:", "", "    quote body", "", "@Since 1.2.0", "", "@Indicator: Server only", "", "@Required", "", "Description."}
	wantOrder := []string{"Description.", "", "@Required", "", "@Indicator: Server only", "", "@Since 1.2.0", "", "@Quote:", "", "    quote body", "", "@Tab \"Example\"", "", "    body"}
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
