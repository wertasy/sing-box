package docgen

import (
	"encoding/json"
	"testing"
)

func TestUnionZeroValue(t *testing.T) {
	array := &SchemaNode{Type: json.RawMessage(`"array"`)}
	integer := &SchemaNode{Type: json.RawMessage(`"integer"`)}
	boolean := &SchemaNode{Type: json.RawMessage(`"boolean"`)}
	str := &SchemaNode{Type: json.RawMessage(`"string"`)}
	object := &SchemaNode{Type: json.RawMessage(`"object"`)}

	cases := []struct {
		name string
		node *SchemaNode
		want string
		ok   bool
	}{
		{"listable", &SchemaNode{AnyOf: []*SchemaNode{array, str}}, "[]", true},
		{"mark", &SchemaNode{AnyOf: []*SchemaNode{str, integer}}, "0", true},
		{"bool flag", &SchemaNode{AnyOf: []*SchemaNode{boolean, str}}, "false", true},
		{"string or object", &SchemaNode{AnyOf: []*SchemaNode{str, object}}, `""`, true},
		{"objectish", &SchemaNode{Properties: &SchemaNodeMap{}}, "", false},
		{"empty", &SchemaNode{}, "", false},
	}
	for _, tc := range cases {
		got, ok := unionZeroValue(tc.node, &SchemaIndex{Defs: map[string]*SchemaNode{}})
		if got != tc.want || ok != tc.ok {
			t.Errorf("%s: unionZeroValue = %q,%v want %q,%v", tc.name, got, ok, tc.want, tc.ok)
		}
	}
}

func TestExampleTag(t *testing.T) {
	ast, err := LoadASTIndex("../option")
	if err != nil {
		t.Fatal(err)
	}
	flat, err := ast.Flatten("SSHOutboundOptions", nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"user":             `"root"`,
		"private_key_path": `"$HOME/.ssh/id_rsa"`,
		// The corpus spells this listable field as an empty string, so the
		// zero value must be explicit; listable marks keep their auto form.
		"private_key": `""`,
	}
	for _, entry := range flat {
		expected, ok := want[entry.Key]
		if !ok {
			continue
		}
		if entry.Tags.Example != expected {
			t.Errorf("%s example = %q, want %q", entry.Key, entry.Tags.Example, expected)
		}
		delete(want, entry.Key)
	}
	if len(want) != 0 {
		t.Errorf("missing example entries: %v", want)
	}
	for _, entry := range flat {
		if entry.Key == "network_type" && entry.Tags.Example != "" {
			t.Errorf("zero example must be omitted, got %q", entry.Tags.Example)
		}
	}
	// Zero values are omitted from comments by rule.
	flat, err = ast.Flatten("MASQUEServerEndpointOptions", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range flat {
		if entry.Key == "mtu" && entry.Tags.Example == "" {
			t.Error("non-zero mtu example missing")
		}
		if entry.Key == "tls" && entry.Tags.Example != "" {
			t.Errorf("zero example must be omitted, got %q", entry.Tags.Example)
		}
	}
}

func TestTagSample(t *testing.T) {
	ast, err := LoadASTIndex("../option")
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]string{
		"SSHOutboundOptions":           "ssh-out",
		"ShadowsocksOutboundOptions":   "ss-out",
		"MASQUEServerEndpointOptions":  "masque-server",
		"OpenVPNClientEndpointOptions": "ovpn-client",
		"TunInboundOptions":            "tun-in",
	}
	for root, want := range cases {
		got, ok := ast.TagSample(root)
		if !ok || got != want {
			t.Errorf("TagSample(%s) = %q,%v want %q", root, got, ok, want)
		}
	}
}

func TestEmbedExampleValues(t *testing.T) {
	ast, err := LoadASTIndex("../option")
	if err != nil {
		t.Fatal(err)
	}
	check := func(root, key, want string) {
		t.Helper()
		flat, err := ast.Flatten(root, nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range flat {
			if entry.Key != key {
				continue
			}
			if entry.Tags.Example != want {
				t.Errorf("%s.%s example = %q, want %q", root, key, entry.Tags.Example, want)
			}
			return
		}
		t.Fatalf("%s.%s not found", root, key)
	}
	check("SSHOutboundOptions", "server", `"127.0.0.1"`)
	check("SSHOutboundOptions", "server_port", "22")
	check("HTTPOutboundOptions", "server", `"127.0.0.1"`)
	check("HTTPOutboundOptions", "server_port", "1080")
	check("_MASQUEClientEndpointOptions", "server_port", "443")
}
