package config

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Service ServerConfig `hcl:"server,block"`
	Log     LogConfig    `hcl:"log,block"`
}

type ServerConfig struct {
	ListenAddr string `hcl:"listen_addr"`
}
type LogConfig struct {
	Format string `hcl:"format"`
	Level  string `hcl:"level"`
}

func ParseConfig(path string) (*Config, error) {
	var config Config

	err := hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
