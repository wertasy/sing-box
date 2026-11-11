package option

import "github.com/sagernet/sing/common/json/badoption"

type SSHOutboundOptions struct {
	DialerOptions
	ServerOptions
	// SSH user, root will be used if empty.
	User string `json:"user,omitempty"`
	// Password.
	Password string `json:"password,omitempty"`
	// Private key.
	PrivateKey badoption.Listable[string] `json:"private_key,omitempty"`
	// Private key path.
	PrivateKeyPath string `json:"private_key_path,omitempty"`
	// Private key passphrase.
	PrivateKeyPassphrase string `json:"private_key_passphrase,omitempty"`
	// Host key. Accept any if empty.
	HostKey badoption.Listable[string] `json:"host_key,omitempty"`
	// Host key algorithms.
	HostKeyAlgorithms LegacyListable[string] `json:"host_key_algorithms,omitempty"`
	// Client version. Random version will be used if empty.
	ClientVersion string `json:"client_version,omitempty"`
	// @since 1.14.0
	//
	// Allowed ciphers. Default values are used if empty.
	Cipher badoption.Listable[string] `json:"cipher,omitempty"`
	// @since 1.14.0
	//
	// Allowed MAC algorithms. Default values are used if empty.
	MAC badoption.Listable[string] `json:"mac,omitempty"`
	// @since 1.14.0
	//
	// Allowed key exchange algorithms. Default values are used if empty.
	KexAlgorithm badoption.Listable[string] `json:"kex_algorithm,omitempty"`
}
