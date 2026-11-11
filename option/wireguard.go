package option

import (
	"net/netip"

	"github.com/sagernet/sing/common/json/badoption"
)

type WireGuardEndpointOptions struct {
	// Use system interface.
	//
	// Requires privilege and cannot conflict with exists system interfaces.
	System bool `json:"system,omitempty"`
	// Custom interface name for system interface.
	Name string `json:"name,omitempty"`
	// WireGuard MTU.
	//
	// `1408` will be used by default.
	MTU uint32 `json:"mtu,omitempty"`
	// @required
	//
	// List of IP (v4 or v6) address prefixes to be assigned to the interface.
	Address badoption.Listable[netip.Prefix] `json:"address"`
	// @required
	//
	// WireGuard requires base64-encoded public and private keys. These can be generated using the wg(8) utility:
	//
	// ```shell
	// wg genkey
	// echo "private key" || wg pubkey
	// ```
	//
	// or `sing-box generate wg-keypair`.
	PrivateKey string `json:"private_key"`
	// Listen port for incoming connections.
	ListenPort uint16 `json:"listen_port,omitempty"`
	// @required
	//
	// List of WireGuard peers.
	Peers        []WireGuardPeer    `json:"peers,omitempty"`
	UDPTimeout   badoption.Duration `json:"udp_timeout,omitempty"`
	UDPMapping   UDPNATBehavior     `json:"udp_mapping,omitempty"`
	UDPFiltering UDPNATBehavior     `json:"udp_filtering,omitempty"`
	UDPNATMax    uint32             `json:"udp_nat_max,omitempty"`
	// WireGuard worker count.
	//
	// CPU count is used by default.
	Workers int `json:"workers,omitempty"`
	// @since 1.15.0
	//
	// Allow the endpoint to be disconnected when necessary.
	OnDemand bool `json:"on_demand,omitempty"`
	DialerOptions
}

type WireGuardPeer struct {
	// WireGuard peer address.
	Address string `json:"address,omitempty"`
	// WireGuard peer port.
	Port uint16 `json:"port,omitempty"`
	// @required
	//
	// WireGuard peer public key.
	PublicKey string `json:"public_key,omitempty"`
	// WireGuard peer pre-shared key.
	PreSharedKey string `json:"pre_shared_key,omitempty"`
	// @required
	//
	// WireGuard allowed IPs.
	AllowedIPs LegacyListable[netip.Prefix] `json:"allowed_ips,omitempty"`
	// WireGuard persistent keepalive interval, in seconds.
	//
	// Disabled by default.
	PersistentKeepaliveInterval uint16 `json:"persistent_keepalive_interval,omitempty"`
	// WireGuard reserved field bytes.
	Reserved []uint8 `json:"reserved,omitempty"`
}
