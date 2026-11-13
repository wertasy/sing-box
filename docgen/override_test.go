package docgen

import (
	"os"
	"path/filepath"
	"testing"
)

func writeProbeFile(t *testing.T, dir string, name string, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Embed overrides declared closer to the documented root must win over the
// ones declared on deeper embeds, and every `@Override` body supplies the
// prose of the flattened field.
func TestEmbedOverrideInheritance(t *testing.T) {
	dir := t.TempDir()
	writeProbeFile(t, dir, "probe.go", `package option

type BaseOptions struct {
	// The base port.
	//
	// `+"`53`"+` will be used by default.
	Port int `+"`json:\"port,omitempty\"`"+`
}

type MidOptions struct {
	// @Override port:
	//    The mid port.
	//
	//    `+"`853`"+` will be used by default.
	//
	BaseOptions
}

type TopOptions struct {
	// @Override port:
	//    The top port.
	//
	//    `+"`443`"+` will be used by default.
	//
	MidOptions
}
`)
	index, err := LoadASTIndex(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, testCase := range []struct {
		root      string
		wantProse string
	}{
		{root: "BaseOptions", wantProse: "The base port."},
		{root: "MidOptions", wantProse: "The mid port."},
		{root: "TopOptions", wantProse: "The top port."},
	} {
		entries, err := index.Flatten(testCase.root, nil)
		if err != nil {
			t.Fatal(err)
		}
		var found bool
		for _, entry := range entries {
			if entry.Key != "port" {
				continue
			}
			found = true
			if len(entry.Prose) == 0 || entry.Prose[0] != testCase.wantProse {
				t.Fatalf("%s: port prose = %v, want first line %q", testCase.root, entry.Prose, testCase.wantProse)
			}
		}
		if !found {
			t.Fatalf("%s: port entry missing", testCase.root)
		}
	}
}

// `@Doc <key>:` blocks on a type's doc comment must be collected for every
// marker, preserving indented body lines without swallowing the next block.
func TestStructDocBlocks(t *testing.T) {
	dir := t.TempDir()
	writeProbeFile(t, dir, "probe.go", `package option

// @Doc first key:
//    First body line.
//
//        Indented body line.
//
// @Doc second:
//    Second body line.
type DocumentedOptions struct {
	Field string `+"`json:\"field,omitempty\"`"+`
}
`)
	index, err := LoadASTIndex(dir)
	if err != nil {
		t.Fatal(err)
	}
	extras := index.StructExtras("DocumentedOptions")
	if len(extras) != 2 {
		t.Fatalf("expected 2 documented keys, got %v", extras)
	}
	first := extras["first key"]
	want := []string{"First body line.", "", "    Indented body line."}
	if len(first) != len(want) {
		t.Fatalf("first key body = %#v, want %#v", first, want)
	}
	for index, line := range want {
		if first[index] != line {
			t.Fatalf("first key body = %#v, want %#v", first, want)
		}
	}
	if second := extras["second"]; len(second) != 1 || second[0] != "Second body line." {
		t.Fatalf("second key body = %#v", second)
	}
}

// `@Doc <key>:` blocks declared on a variant holder field must register under
// the enclosing struct together with the holder whose documented fields
// accompany the block.
func TestFieldDocBlocks(t *testing.T) {
	dir := t.TempDir()
	writeProbeFile(t, dir, "probe.go", `package option

type HolderOptions struct {
	// @Doc Guide:
	//    Guide body line.
	Inner InnerOptions `+"`json:\"-\"`"+`
}

type InlineOptions struct {
	InnerOptions `+"`json:\"inline\"`"+`
}

type InnerOptions struct {
	Plain string `+"`json:\"plain,omitempty\"`"+`
	// Documented value.
	//
	// @Since 1.13.0
	Value string `+"`json:\"value,omitempty\"`"+`
}
`)
	index, err := LoadASTIndex(dir)
	if err != nil {
		t.Fatal(err)
	}
	extras := index.StructExtras("HolderOptions")
	if len(extras) != 1 || len(extras["Guide"]) == 0 {
		t.Fatalf("holder doc block missing: %#v", extras)
	}
	if holder := index.StructDocHolder("HolderOptions", "Guide"); holder != "InnerOptions" {
		t.Fatalf("holder = %q", holder)
	}
	fields := index.DocumentedFields("InnerOptions")
	if len(fields) != 1 || fields[0].Key != "value" {
		t.Fatalf("documented fields = %#v", fields)
	}
}

func TestDocPlacementValidation(t *testing.T) {
	for _, testCase := range []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name: "type-level key backed by a field",
			source: `package option

// @Doc port:
//    Wrong place.
type BadOptions struct {
	Port int ` + "`json:\"port,omitempty\"`" + `
}
`,
			wantErr: true,
		},
		{
			name: "field-level key unrelated to its field",
			source: `package option

type BadOptions struct {
	// @Doc other:
	//    Wrong place.
	Port int ` + "`json:\"port,omitempty\"`" + `
}
`,
			wantErr: true,
		},
		{
			name: "field-level nested key under its field",
			source: `package option

type GoodOptions struct {
	// @Doc users.username:
	//    Username.
	Users []string ` + "`json:\"users,omitempty\"`" + `
}
`,
			wantErr: false,
		},
		{
			name: "holder field key",
			source: `package option

type GoodOptions struct {
	// @Doc Guide:
	//    Guide body.
	Inner InnerOptions ` + "`json:\"-\"`" + `
}

type InnerOptions struct {
	Value string ` + "`json:\"value,omitempty\"`" + `
}
`,
			wantErr: false,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			dir := t.TempDir()
			writeProbeFile(t, dir, "probe.go", testCase.source)
			_, err := LoadASTIndex(dir)
			if testCase.wantErr && err == nil {
				t.Fatal("expected placement error")
			}
			if !testCase.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
