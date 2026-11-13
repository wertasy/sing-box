package docgen

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// The reference-only action option types must keep the section-specific
// wording; the runtime types stay shared.
func TestDocsOnlyActionOptions(t *testing.T) {
	index, err := LoadASTIndex("../option")
	if err != nil {
		t.Fatal(err)
	}
	for _, testCase := range []struct {
		root string
		key  string
		want string
	}{
		{root: "BypassActionOptions", key: "outbound", want: "pre-match"},
		{root: "DNSRejectActionOptions", key: "method", want: "REFUSED"},
	} {
		entries, err := index.Flatten(testCase.root, nil)
		if err != nil {
			t.Fatal(err)
		}
		var found bool
		for _, entry := range entries {
			if entry.Key != testCase.key {
				continue
			}
			found = true
			if !strings.Contains(strings.Join(entry.Prose, "\n"), testCase.want) {
				t.Fatalf("%s.%s prose = %q, want %q", testCase.root, testCase.key, entry.Prose, testCase.want)
			}
		}
		if !found {
			t.Fatalf("%s.%s entry missing", testCase.root, testCase.key)
		}
	}
}

// The ACME DNS01 challenge schema is a provider union; every variant must
// expose the provider-specific fields documented on the reference page.
func TestACMESchemaProviderVariants(t *testing.T) {
	content, err := os.ReadFile("../docs/schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var document struct {
		Defs map[string]struct {
			OneOf []struct {
				Properties map[string]json.RawMessage `json:"properties"`
			} `json:"oneOf"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal(content, &document); err != nil {
		t.Fatal(err)
	}
	union, loaded := document.Defs["ACMEProviderDNS01Challenge"]
	if !loaded || len(union.OneOf) != 3 {
		t.Fatalf("ACMEProviderDNS01Challenge union = %+v", union)
	}
	expected := [][]string{
		{"access_key_id", "access_key_secret", "region_id", "security_token"},
		{"api_token", "zone_token"},
		{"username", "password", "subdomain", "server_url"},
	}
	for index, variant := range union.OneOf {
		for _, field := range expected[index] {
			if _, found := variant.Properties[field]; !found {
				t.Fatalf("union variant %d is missing field %q", index, field)
			}
		}
	}
}

// Pages under registry directories get an implicit definition instead of a
// placeholder `defs: [{}]` entry in the manifest.
func TestInferImplicitDefs(t *testing.T) {
	pages := []Page{
		{File: "inbound/plain.md"},
		{File: "inbound/plain.zh.md"},
		{File: "shared/plain.md"},
		{File: "inbound/explicit.md", Defs: []DefRef{{Name: "explicit"}}},
		{File: "inbound/sectioned.md", Sections: []Section{{Heading: "### Fields"}}},
	}
	inferImplicitDefs(pages)
	if len(pages[0].Defs) != 1 {
		t.Fatalf("registry page defs = %#v", pages[0].Defs)
	}
	if def := pages[0].Defs[0]; def.Name != "" || def.Root != "" || def.Registry != "" {
		t.Fatalf("implicit def = %#v", def)
	}
	for _, index := range []int{1, 2, 4} {
		if len(pages[index].Defs) != 0 {
			t.Fatalf("%s: unexpected inferred defs %#v", pages[index].File, pages[index].Defs)
		}
	}
	if len(pages[3].Defs) != 1 || pages[3].Defs[0].Name != "explicit" {
		t.Fatalf("explicit defs changed: %#v", pages[3].Defs)
	}
}
