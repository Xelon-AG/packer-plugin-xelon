package xelonapi

import (
	"github.com/hashicorp/packer-plugin-sdk/useragent"

	"github.com/Xelon-AG/packer-plugin-xelon/internal/config"
	"github.com/Xelon-AG/packer-plugin-xelon/internal/version"
	"github.com/Xelon-AG/xelon-sdk-go/xelon"
)

func NewXelonClient(c config.AccessConfig) *xelon.Client {
	opts := []xelon.ClientOption{xelon.WithUserAgent(useragent.String(version.PluginVersion.FormattedVersion()))}
	if c.BaseURL != "" {
		opts = append(opts, xelon.WithBaseURL(c.BaseURL))
	}
	opts = append(opts, xelon.WithClientID(c.ClientID))

	return xelon.NewClient(c.Token, opts...)
}
