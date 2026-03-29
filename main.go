package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"

	xelonBuilder "github.com/Xelon-AG/packer-plugin-xelon/internal/builder"
	xelonVersion "github.com/Xelon-AG/packer-plugin-xelon/internal/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder(plugin.DEFAULT_NAME, new(xelonBuilder.Builder))
	pps.SetVersion(xelonVersion.PluginVersion)
	err := pps.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
