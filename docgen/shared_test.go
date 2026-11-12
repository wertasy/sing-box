package docgen

import "testing"

func TestSharedOrigins(t *testing.T) {
	ast, err := LoadASTIndex("../option")
	if err != nil {
		t.Fatal(err)
	}
	collect := func(root string) (map[string]string, []string) {
		t.Helper()
		origins, err := ast.SharedOrigins(root)
		if err != nil {
			t.Fatalf("SharedOrigins(%s): %v", root, err)
		}
		byOrigin := make(map[string]string, len(origins))
		var order []string
		for _, origin := range origins {
			byOrigin[origin.Origin] = origin.Slug
			order = append(order, origin.Origin)
		}
		return byOrigin, order
	}

	t.Run("dialer embeds", func(t *testing.T) {
		byOrigin, _ := collect("SSHOutboundOptions")
		if byOrigin["DialerOptions"] != "dial" {
			t.Errorf("DialerOptions slug = %q", byOrigin["DialerOptions"])
		}
		if byOrigin["AbstractDialerOptions"] != "dial" {
			t.Errorf("AbstractDialerOptions slug = %q", byOrigin["AbstractDialerOptions"])
		}
		if _, found := byOrigin["ServerOptions"]; found {
			t.Errorf("ServerOptions must not be shared: %v", byOrigin)
		}
	})

	t.Run("listener embeds", func(t *testing.T) {
		byOrigin, _ := collect("RedirectInboundOptions")
		if byOrigin["ListenOptions"] != "listen" {
			t.Errorf("ListenOptions slug = %q", byOrigin["ListenOptions"])
		}
		if byOrigin["InboundOptions"] != "listen" {
			t.Errorf("InboundOptions slug = %q", byOrigin["InboundOptions"])
		}
	})

	t.Run("no shared embeds", func(t *testing.T) {
		byOrigin, _ := collect("CacheFileOptions")
		if len(byOrigin) != 0 {
			t.Errorf("CacheFileOptions shared origins = %v", byOrigin)
		}
	})

	t.Run("root itself is not shared", func(t *testing.T) {
		byOrigin, _ := collect("DialerOptions")
		if len(byOrigin) != 0 {
			t.Errorf("DialerOptions root shared origins = %v", byOrigin)
		}
	})
}
