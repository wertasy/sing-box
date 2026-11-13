package option

import (
	"encoding/json"

	"github.com/sagernet/sing-box/schema"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json/badjson"
)

// @Tag st-in
type ShadowTLSInboundOptions struct {
	ListenOptions

	// ShadowTLS protocol version.
	//
	// | Value         | Protocol Version                                                                        |
	// |---------------|-----------------------------------------------------------------------------------------|
	// | `1` (default) | [ShadowTLS v1](https://github.com/ihciah/shadow-tls/blob/master/docs/protocol-en.md#v1) |
	// | `2`           | [ShadowTLS v2](https://github.com/ihciah/shadow-tls/blob/master/docs/protocol-en.md#v2) |
	// | `3`           | [ShadowTLS v3](https://github.com/ihciah/shadow-tls/blob/master/docs/protocol-v3-en.md) |
	//
	// @Example: 3
	Version int `json:"version,omitempty" enum:"1,2,3"`

	// ShadowTLS password.
	//
	// Only available in the ShadowTLS protocol 2.
	//
	// @Example: "fuck me till the daylight"
	Password string `json:"password,omitempty"`

	// ShadowTLS users.
	//
	// Only available in the ShadowTLS protocol 3.
	//
	// @Example:
	//    [
	//      {
	//        "name": "sekai",
	//        "password": "8JCsPssfgS8tiRwiMlhARg=="
	//      }
	//    ]
	Users []ShadowTLSUser `json:"users,omitempty"`

	// When `wildcard_sni` is configured to `all`, the server address is optional.
	//
	// Handshake server address and [Dial Fields](/configuration/shared/dial/).
	//
	// @Example:
	//    {
	//      "server": "google.com",
	//      "server_port": 443,
	//      ... // Dial Fields
	//    }
	//
	//      ... // Dial Fields
	//    }
	//
	// @Required
	Handshake ShadowTLSHandshakeOptions `json:"handshake,omitempty"`

	// Handshake server address and [Dial Fields](/configuration/shared/dial/) for specific server name.
	//
	// Only available in the ShadowTLS protocol 2/3.
	//
	// @Example:
	//    {
	//      "example.com": {
	//        "server": "example.com",
	//        "server_port": 443,
	//        ... // Dial Fields
	//      }
	//    }
	//
	//        ... // Dial Fields
	//      }
	//    }
	HandshakeForServerName *badjson.TypedMap[string, ShadowTLSHandshakeOptions] `json:"handshake_for_server_name,omitempty"`

	// ShadowTLS strict mode.
	//
	// Only available in the ShadowTLS protocol 3.
	//
	// @Example: false
	StrictMode bool `json:"strict_mode,omitempty"`

	// ShadowTLS wildcard SNI mode.
	//
	// Available values are:
	//
	// * `off`: (default) Disabled.
	// * `authed`: Authenticated connections will have their destination overwritten to `(servername):443`
	// * `all`: All connections will have their destination overwritten to `(servername):443`
	//
	// Additionally, connections matching `handshake_for_server_name` are not affected.
	//
	// Only available in the ShadowTLS protocol 3.
	//
	// @Since 1.12.0
	//
	// @Example: ""
	WildcardSNI WildcardSNI `json:"wildcard_sni,omitempty"`
}

type WildcardSNI int

const (
	ShadowTLSWildcardSNIOff WildcardSNI = iota
	ShadowTLSWildcardSNIAuthed
	ShadowTLSWildcardSNIAll
)

func (w WildcardSNI) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

func (w WildcardSNI) String() string {
	switch w {
	case ShadowTLSWildcardSNIOff:
		return "off"
	case ShadowTLSWildcardSNIAuthed:
		return "authed"
	case ShadowTLSWildcardSNIAll:
		return "all"
	default:
		panic("unknown wildcard SNI value")
	}
}

func (w *WildcardSNI) UnmarshalJSON(bytes []byte) error {
	var valueString string
	err := json.Unmarshal(bytes, &valueString)
	if err != nil {
		return err
	}
	switch valueString {
	case "off", "":
		*w = ShadowTLSWildcardSNIOff
	case "authed":
		*w = ShadowTLSWildcardSNIAuthed
	case "all":
		*w = ShadowTLSWildcardSNIAll
	default:
		return E.New("unknown wildcard SNI value: ", valueString)
	}
	return nil
}

func (w WildcardSNI) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.StringEnum("", "off", "authed", "all"), nil
}

type ShadowTLSUser struct {
	// User name.
	Name string `json:"name,omitempty"`

	// User password.
	Password string `json:"password,omitempty"`
}

type ShadowTLSHandshakeOptions struct {
	ServerOptions
	DialerOptions
}

// @Tag st-out
type ShadowTLSOutboundOptions struct {
	DialerOptions

	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 1080
	ServerOptions

	// ShadowTLS protocol version.
	//
	// | Value         | Protocol Version                                                                        |
	// |---------------|-----------------------------------------------------------------------------------------|
	// | `1` (default) | [ShadowTLS v1](https://github.com/ihciah/shadow-tls/blob/master/docs/protocol-en.md#v1) |
	// | `2`           | [ShadowTLS v2](https://github.com/ihciah/shadow-tls/blob/master/docs/protocol-en.md#v2) |
	// | `3`           | [ShadowTLS v3](https://github.com/ihciah/shadow-tls/blob/master/docs/protocol-v3-en.md) |
	//
	// @Example: 3
	Version int `json:"version,omitempty" enum:"1,2,3"`

	// Set password.
	//
	// Only available in the ShadowTLS v2/v3 protocol.
	//
	// @Example: "fuck me till the daylight"
	Password string `json:"password,omitempty"`

	// @Override:
	//    @Required
	//
	//    TLS configuration, see [TLS](/configuration/shared/tls/#outbound).
	OutboundTLSOptionsContainer
}
