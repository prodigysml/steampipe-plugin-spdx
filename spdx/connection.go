package spdx

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type spdxConfig struct{}

var ConfigSchema = map[string]*schema.Attribute{}

func ConfigInstance() interface{} {
	return &spdxConfig{}
}

func GetConfig(connection *plugin.Connection) spdxConfig {
	if connection == nil || connection.Config == nil {
		return spdxConfig{}
	}
	config, _ := connection.Config.(spdxConfig)
	return config
}
