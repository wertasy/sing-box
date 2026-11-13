package option

import (
	"net/netip"

	"github.com/sagernet/sing/common/json/badoption"
)

// @Tag wg-ep
type WireGuardEndpointOptions struct {
	// Use system interface.
	//
	// Requires privilege and cannot conflict with exists system interfaces.
	//
	// @Example: false
	System bool `json:"system,omitempty"`

	// Custom interface name for system interface.
	//
	// @Example: ""
	Name string `json:"name,omitempty"`

	// WireGuard MTU.
	//
	// @Default: 1408
	//
	// @Example: 1408
	MTU uint32 `json:"mtu,omitempty"`

	// List of IP (v4 or v6) address prefixes to be assigned to the interface.
	//
	// @Required
	//
	// @Example: []
	Address badoption.Listable[netip.Prefix] `json:"address"`

	// WireGuard requires base64-encoded public and private keys. These can be generated using the wg(8) utility:
	//
	// ```shell
	// wg genkey
	// echo "private key" || wg pubkey
	// ```
	//
	// or `sing-box generate wg-keypair`.
	//
	// @Required
	//
	// @Example: ""
	PrivateKey string `json:"private_key"`

	// Listen port for incoming connections.
	//
	// @Example: 10000
	ListenPort uint16 `json:"listen_port,omitempty"`

	// List of WireGuard peers.
	//
	// @Example:
	//    [
	//      {
	//        "address": "127.0.0.1",
	//        "port": 10001,
	//        "public_key": "",
	//        "pre_shared_key": "",
	//        "allowed_ips": [],
	//        "persistent_keepalive_interval": 0,
	//        "reserved": [0, 0, 0]
	//      }
	//    ]
	//
	// @Required
	Peers []WireGuardPeer `json:"peers,omitempty"`

	// @Omit
	//
	// @Example: ""
	UDPTimeout badoption.Duration `json:"udp_timeout,omitempty"`

	// @Omit
	//
	// @Example: ""
	UDPMapping UDPNATBehavior `json:"udp_mapping,omitempty"`

	// @Omit
	//
	// @Example: ""
	UDPFiltering UDPNATBehavior `json:"udp_filtering,omitempty"`

	// @Omit
	//
	// @Example: 0
	UDPNATMax uint32 `json:"udp_nat_max,omitempty"`

	// WireGuard worker count.
	//
	// CPU count will be used by default.
	//
	// @Example: 0
	Workers int `json:"workers,omitempty"`

	// Allow the endpoint to be disconnected when necessary.
	//
	// @Since 1.15.0
	//
	// @Example: false
	OnDemand bool `json:"on_demand,omitempty"`
	DialerOptions
}

type WireGuardPeer struct {
	// WireGuard peer address.
	Address string `json:"address,omitempty"`

	// WireGuard peer port.
	Port uint16 `json:"port,omitempty"`

	// WireGuard peer public key.
	//
	// @Required
	PublicKey string `json:"public_key,omitempty"`

	// WireGuard peer pre-shared key.
	PreSharedKey string `json:"pre_shared_key,omitempty"`

	// WireGuard allowed IPs.
	//
	// @Required
	AllowedIPs LegacyListable[netip.Prefix] `json:"allowed_ips,omitempty"`

	// WireGuard persistent keepalive interval, in seconds.
	//
	// Disabled by default.
	PersistentKeepaliveInterval uint16 `json:"persistent_keepalive_interval,omitempty"`

	// WireGuard reserved field bytes.
	Reserved []uint8 `json:"reserved,omitempty"`
}
