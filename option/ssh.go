package option

import "github.com/sagernet/sing/common/json/badoption"

// @Tag ssh-out
type SSHOutboundOptions struct {
	DialerOptions

	// @Override server:
	//    @Required
	//
	//    Server address.
	//
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    Server port. 22 will be used if empty.
	//
	//    @Example: 22
	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 22
	ServerOptions

	// SSH user, root will be used if empty.
	//
	// @Example: "root"
	User string `json:"user,omitempty"`

	// Password.
	//
	// @Example: "admin"
	Password string `json:"password,omitempty"`

	// Private key.
	//
	// @Example: ""
	PrivateKey badoption.Listable[string] `json:"private_key,omitempty"`

	// Private key path.
	//
	// @Example: "$HOME/.ssh/id_rsa"
	PrivateKeyPath string `json:"private_key_path,omitempty"`

	// Private key passphrase.
	//
	// @Example: ""
	PrivateKeyPassphrase string `json:"private_key_passphrase,omitempty"`

	// Host key. Accept any if empty.
	//
	// @Example:
	//    [
	//      "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdH..."
	//    ]
	HostKey badoption.Listable[string] `json:"host_key,omitempty"`

	// Host key algorithms.
	//
	// @Example: []
	HostKeyAlgorithms LegacyListable[string] `json:"host_key_algorithms,omitempty"`

	// Client version. Random version will be used if empty.
	//
	// @Example: "SSH-2.0-OpenSSH_7.4p1"
	ClientVersion string `json:"client_version,omitempty"`

	// Allowed ciphers. Default values are used if empty.
	//
	// @Since 1.14.0
	//
	// @Example: []
	Cipher badoption.Listable[string] `json:"cipher,omitempty"`

	// Allowed MAC algorithms. Default values are used if empty.
	//
	// @Since 1.14.0
	//
	// @Example: []
	MAC badoption.Listable[string] `json:"mac,omitempty"`

	// Allowed key exchange algorithms. Default values are used if empty.
	//
	// @Since 1.14.0
	//
	// @Example: []
	KexAlgorithm badoption.Listable[string] `json:"kex_algorithm,omitempty"`
}
