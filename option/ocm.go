package option

import (
	"github.com/sagernet/sing/common/json/badoption"
)

type OCMServiceOptions struct {
	ListenOptions
	InboundTLSOptionsContainer

	// Path to the OpenAI OAuth credentials file.
	//
	// If not specified, defaults to:
	// - `$CODEX_HOME/auth.json` if `CODEX_HOME` environment variable is set
	// - `~/.codex/auth.json` otherwise
	//
	// Refreshed tokens are automatically written back to the same location.
	CredentialPath string `json:"credential_path,omitempty"`

	// List of authorized users for token authentication.
	//
	// If empty, no authentication is required.
	//
	// Object format:
	//
	// ```json
	// {
	//   "name": "",
	//   "token": ""
	// }
	// ```
	//
	// Object fields:
	//
	// - `name`: Username identifier for tracking purposes.
	// - `token`: Bearer token for authentication. Clients authenticate by setting the `Authorization: Bearer <token>` header.
	Users []OCMUser `json:"users,omitempty"`

	// Custom HTTP headers to send to the OpenAI API.
	//
	// These headers will override any existing headers with the same name.
	Headers badoption.HTTPHeader `json:"headers,omitempty"`

	// Outbound tag for connecting to the OpenAI API.
	Detour string `json:"detour,omitempty" reference:"outbound"`

	// Path to the file for storing aggregated API usage statistics.
	//
	// Usage tracking is disabled if not specified.
	//
	// When enabled, the service tracks and saves comprehensive statistics including:
	// - Request counts
	// - Token usage (input, output, cached)
	// - Calculated costs in USD based on OpenAI API pricing
	//
	// Statistics are organized by model and optionally by user when authentication is enabled.
	//
	// The statistics file is automatically saved every minute and upon service shutdown.
	UsagesPath string `json:"usages_path,omitempty"`
}

type OCMUser struct {
	Name  string `json:"name,omitempty"`
	Token string `json:"token,omitempty"`
}
