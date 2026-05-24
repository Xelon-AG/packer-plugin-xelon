//go:generate packer-sdc struct-markdown
package internal

import (
	"errors"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

var (
	ErrMissingClientID = errors.New(`client id is required: set "client_id" attribute or XELON_CLIENT_ID environment variable`)
	ErrMissingToken    = errors.New(`token is required: set "token" attribute or XELON_TOKEN environment variable`)
)

// AccessConfig is for common configuration related to Xelon HQ API access.
type AccessConfig struct {
	// The base URL endpoint for Xelon HQ. Default is `https://hq.xelon.ch/api/v2/`.
	// Alternatively, can be configured using the `XELON_BASE_URL` environment variable.
	BaseURL string `mapstructure:"base_url" required:"false"`
	// The client ID for IP ranges.
	// Alternatively, can be configured using the `XELON_CLIENT_ID` environment variable.
	ClientID string `mapstructure:"client_id" required:"true"`
	// The Xelon access token.
	// Alternatively, can be configured using the `XELON_TOKEN` environment variable.
	Token string `mapstructure:"token" required:"true"`
}

func (c *AccessConfig) Prepare(_ *interpolate.Context, _ ...any) *packer.MultiError {
	var errs *packer.MultiError

	if c.BaseURL == "" {
		if os.Getenv("XELON_BASE_URL") != "" {
			c.BaseURL = os.Getenv("XELON_BASE_URL")
		}
	}

	if c.ClientID == "" {
		c.ClientID = os.Getenv("XELON_CLIENT_ID")
	}
	if c.ClientID == "" {
		errs = packer.MultiErrorAppend(errs, ErrMissingClientID)
	}

	if c.Token == "" {
		c.Token = os.Getenv("XELON_TOKEN")
	}
	if c.Token == "" {
		errs = packer.MultiErrorAppend(errs, ErrMissingToken)
	}

	return errs
}
