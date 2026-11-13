package docgen

import (
	"slices"
	"strings"
	"testing"
)

func TestOmitPages(t *testing.T) {
	// Page-scoped omit semantics (pure function).
	unscoped := FieldTags{Omit: true}
	if !omittedOnPage(unscoped, "inbound/a.md") || !omittedOnPage(unscoped, "inbound/a.zh.md") {
		t.Error("bare @Omit must hide on every page")
	}
	scoped := FieldTags{Omit: true, OmitPages: []string{"inbound/a.md"}}
	if !omittedOnPage(scoped, "inbound/a.md") || !omittedOnPage(scoped, "inbound/a.zh.md") {
		t.Error("listed page and its twin must hide")
	}
	if omittedOnPage(scoped, "inbound/b.md") {
		t.Error("unlisted page must keep the field")
	}
	if omittedOnPage(FieldTags{}, "inbound/a.md") {
		t.Error("unmarked field must not hide")
	}

	// Real data: the ECH container hides `enabled` on all pages, so the
	// marker is bare (no page list).
	ast, err := LoadASTIndex("../option")
	if err != nil {
		t.Fatal(err)
	}
	flat, err := ast.Flatten("InboundECHOptions", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range flat {
		if entry.Key != "enabled" {
			continue
		}
		if !entry.Tags.Omit {
			t.Errorf("enabled not marked @Omit: %+v", entry.Tags)
		}
		if len(entry.Tags.OmitPages) != 0 {
			t.Errorf("enabled should use bare @Omit, got pages %v", entry.Tags.OmitPages)
		}
		return
	}
	t.Fatal("enabled not found")
}

func TestRequiredSuffix(t *testing.T) {
	prose := []string{"@Required: with the `default` provider."}
	en := expandProse(prose, "en", "en")
	if len(en) != 1 || en[0] != "==Required== with the `default` provider." {
		t.Errorf("en = %v", en)
	}
	zh := expandProse(prose, "zh", "zh")
	if len(zh) != 1 || zh[0] != "==必填== with the `default` provider." {
		t.Errorf("zh = %v", zh)
	}
}

func TestEmbedOverrides(t *testing.T) {
	ast, err := LoadASTIndex("../option")
	if err != nil {
		t.Fatal(err)
	}
	proseOf := func(root, key string) []string {
		t.Helper()
		flat, err := ast.Flatten(root, nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range flat {
			if entry.Key == key {
				return entry.Prose
			}
		}
		t.Fatalf("%s.%s not found", root, key)
		return nil
	}
	if got := proseOf("SSHOutboundOptions", "server"); len(got) == 0 || got[0] != "@Required" {
		t.Errorf("ssh server override = %v", got)
	}
	if got := proseOf("SSHOutboundOptions", "server_port"); len(got) == 0 || got[0] != "Server port. 22 will be used if empty." {
		t.Errorf("ssh server_port override = %v", got)
	}
	if got := proseOf("Hysteria2InboundOptions", "tls"); len(got) < 3 || got[2] != "TLS configuration, see [TLS](/configuration/shared/tls/#inbound)." {
		t.Errorf("hysteria2 tls override = %v", got)
	}
	realmTLS := strings.Join(proseOf("HysteriaRealmServiceOptions", "tls"), "\n")
	if !strings.Contains(realmTLS, "When configured, the realm serves HTTP/2 over TLS; otherwise plain HTTP/1.1.") {
		t.Errorf("realm tls override = %v", proseOf("HysteriaRealmServiceOptions", "tls"))
	}
	if strings.Contains(realmTLS, "@Required") || strings.Contains(realmTLS, "==Required==") {
		t.Errorf("realm tls override must not require TLS: %v", proseOf("HysteriaRealmServiceOptions", "tls"))
	}
	// Expand path must honor embed overrides too (tun platform.http_proxy).
	flat, err := ast.Flatten("TunInboundOptions", nil)
	if err != nil {
		t.Fatal(err)
	}
	expanded, err := ast.Expand(flat, "platform")
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, entry := range expanded {
		if entry.Key == "platform.http_proxy.server" {
			found = true
			if len(entry.Prose) == 0 || entry.Prose[0] != "@Required" {
				t.Errorf("tun platform server override = %v", entry.Prose)
			}
		}
	}
	if !found {
		t.Error("platform.http_proxy.server not expanded")
	}
	// Unannotated embeds keep their default prose; comments are authored with
	// the description first and the markers below it.
	if got := proseOf("VLESSOutboundOptions", "server"); len(got) < 3 || got[0] != "The server address." || got[len(got)-1] != "@Required" {
		t.Errorf("vless server default = %v", got)
	}
	if got := proseOf("OpenVPNClientEndpointOptions", "system"); len(got) == 0 || got[0] != "Use a system interface." {
		t.Errorf("openvpn-client system default = %v", got)
	}
	if got := proseOf("OpenVPNServerEndpointOptions", "system"); len(got) == 0 || got[0] != "Use system interface." {
		t.Errorf("openvpn-server system override = %v", got)
	}
}

func TestOrderMarkers(t *testing.T) {
	// @Since and @Deprecated in keep their rendered positions above the
	// description (@Since first); the description follows, then the remaining
	// groups in source order.
	prose := []string{"Do not reserve a port.", "", "@Quote:", "", "    Only supported on Linux.", "", "@Since 1.13.0", "", "@Deprecated in 1.14.0:", "", "    Use the replacement."}
	want := []string{"@Since 1.13.0", "", "@Deprecated in 1.14.0:", "", "    Use the replacement.", "", "@Quote:", "", "    Only supported on Linux.", "", "Do not reserve a port."}
	if got := orderMarkers(prose); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("orderMarkers = %q, want %q", got, want)
	}
	// Required banners are "rest" groups: they stay above the description.
	if got := orderMarkers([]string{"@Required", "", "Description."}); strings.Join(got, "\n") != strings.Join([]string{"@Required", "", "Description."}, "\n") {
		t.Fatalf("required banner not kept above description: %q", got)
	}
	if got := orderMarkers([]string{"@Since 1.0.0"}); strings.Join(got, "\n") != "@Since 1.0.0" {
		t.Fatalf("marker-only prose changed: %q", got)
	}
	if got := orderMarkers([]string{"Plain description."}); strings.Join(got, "\n") != "Plain description." {
		t.Fatalf("description-only prose changed: %q", got)
	}
	// Round trip: schema order then render order restores the authored split.
	if got := orderMarkers(reorderDescription(prose)); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("round trip = %q, want %q", got, want)
	}
}

func TestDeprecatedNote(t *testing.T) {
	ast, err := LoadASTIndex("../option")
	if err != nil {
		t.Fatal(err)
	}
	flat, err := ast.Flatten("CacheFileOptions", nil)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, entry := range flat {
		if entry.Key != "store_rdrc" {
			continue
		}
		found = true
		if !entry.Tags.Deprecated {
			t.Error("store_rdrc not marked deprecated")
		}
		// The plain godoc `Deprecated:` line renders the `- Deprecated.` tag;
		// the replacement note is not retained.
		if entry.Tags.DeprecatedNote != "" {
			t.Errorf("godoc Deprecated line must not keep a note, got %q", entry.Tags.DeprecatedNote)
		}
		tagLines, err := renderTagLine([]FlattenEntry{entry}, nil, nil, "en")
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Contains(tagLines, "- Deprecated.") {
			t.Errorf("deprecated tag line missing: %v", tagLines)
		}
	}
	if !found {
		t.Fatal("store_rdrc not found")
	}
}
