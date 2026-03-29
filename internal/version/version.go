package version

import "github.com/hashicorp/packer-plugin-sdk/version"

var (
	Version          = "0.0.0" // 0.0.0 means dev version
	PrereleaseSuffix = "dev"
	Metadata         = ""
	PluginVersion    = version.NewPluginVersion(Version, PrereleaseSuffix, Metadata)
)
