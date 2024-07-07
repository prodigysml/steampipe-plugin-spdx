package main

import (
	"github.com/prodigysml/steampipe-plugin-spdx/spdx"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: spdx.Plugin,
	})
}
