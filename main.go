package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"

	"github.com/Xelon-AG/packer-plugin-xelon/internal/builder"
	"github.com/Xelon-AG/packer-plugin-xelon/internal/datasource"
	"github.com/Xelon-AG/packer-plugin-xelon/internal/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder(plugin.DEFAULT_NAME, new(builder.Builder))
	pps.RegisterDatasource("network", new(datasource.NetworkDatasource))
	pps.RegisterDatasource("template", new(datasource.TemplateDatasource))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
