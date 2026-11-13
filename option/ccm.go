package option

import (
	"github.com/sagernet/sing/common/json/badoption"
)

type CCMServiceOptions struct {
	ListenOptions
	InboundTLSOptionsContainer

	// Path to the Claude Code OAuth credentials file.
	//
	// If not specified, defaults to:
	// - `$CLAUDE_CONFIG_DIR/.credentials.json` if `CLAUDE_CONFIG_DIR` environment variable is set
	// - `~/.claude/.credentials.json` otherwise
	//
	// On macOS, credentials are read from the system keychain first, then fall back to the file if unavailable.
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
	// - `token`: Bearer token for authentication. Claude Code authenticates by setting the `ANTHROPIC_AUTH_TOKEN` environment variable to their token value.
	Users []CCMUser `json:"users,omitempty"`

	// Custom HTTP headers to send to the Claude API.
	//
	// These headers will override any existing headers with the same name.
	Headers badoption.HTTPHeader `json:"headers,omitempty"`

	// Outbound tag for connecting to the Claude API.
	Detour string `json:"detour,omitempty" reference:"outbound"`

	// Path to the file for storing aggregated API usage statistics.
	//
	// Usage tracking is disabled if not specified.
	//
	// When enabled, the service tracks and saves comprehensive statistics including:
	// - Request counts
	// - Token usage (input, output, cache read, cache creation)
	// - Calculated costs in USD based on Claude API pricing
	//
	// Statistics are organized by model, context window (200k standard vs 1M premium), and optionally by user when authentication is enabled.
	//
	// The statistics file is automatically saved every minute and upon service shutdown.
	UsagesPath string `json:"usages_path,omitempty"`
}

type CCMUser struct {
	Name  string `json:"name,omitempty"`
	Token string `json:"token,omitempty"`
}
