package option

import "github.com/sagernet/sing/common/json/badoption"

type AnyTLSInboundOptions struct {
	ListenOptions
	InboundTLSOptionsContainer

	// AnyTLS users.
	//
	// @Example:
	//    [
	//      {
	//        "name": "sekai",
	//        "password": "8JCsPssfgS8tiRwiMlhARg=="
	//      }
	//    ]
	//
	// @Required
	Users []AnyTLSUser `json:"users,omitempty"`

	// AnyTLS padding scheme line array.
	//
	// Default padding scheme:
	//
	// ```json
	// [
	//   "stop=8",
	//   "0=30-30",
	//   "1=100-400",
	//   "2=400-500,c,500-1000,c,500-1000,c,500-1000,c,500-1000",
	//   "3=9-9,500-1000",
	//   "4=500-1000",
	//   "5=500-1000",
	//   "6=500-1000",
	//   "7=500-1000"
	// ]
	// ```
	PaddingScheme badoption.Listable[string] `json:"padding_scheme,omitempty"`
}

type AnyTLSUser struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type AnyTLSOutboundOptions struct {
	DialerOptions

	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 1080
	ServerOptions

	// @Override:
	//    @Required
	//
	//    TLS configuration, see [TLS](/configuration/shared/tls/#outbound).
	OutboundTLSOptionsContainer

	// The AnyTLS password.
	//
	// @Example: "8JCsPssfgS8tiRwiMlhARg=="
	//
	// @Required
	Password string `json:"password,omitempty"`

	// Interval checking for idle sessions. Default: 30s.
	//
	// @Example: "30s"
	IdleSessionCheckInterval badoption.Duration `json:"idle_session_check_interval,omitempty"`

	// In the check, close sessions that have been idle for longer than this. Default: 30s.
	//
	// @Example: "30s"
	IdleSessionTimeout badoption.Duration `json:"idle_session_timeout,omitempty"`

	// In the check, at least the first `n` idle sessions are kept open. Default value: `n`=0
	//
	// @Example: 5
	MinIdleSession int `json:"min_idle_session,omitempty"`

	// @Since 1.13.16
	//
	//     Check [AnyTLS client metadata](/manual/misc/anytls-client-metadata/).
	//
	// @Example: ""
	ClientMetadata string `json:"client_metadata,omitempty"`
}
