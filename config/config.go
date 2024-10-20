package config

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Server ServerConfig `hcl:"server,block"`
	Log    LogConfig    `hcl:"log,block"`
	Nvim   NvimConfig   `hcl:"nvim,block"`
}

type ServerConfig struct {
	ListenAddr string `hcl:"listen_addr"`
}
type LogConfig struct {
	Format string `hcl:"format"`
	Level  string `hcl:"level"`
}

type NvimConfig struct {
	Cmd            string   `hcl:"cmd"`
	Args           []string `hcl:"args"`
	ForwardEnvVars []string `hcl:"forwardEnvVars"`
}

func ParseConfig(path string) (*Config, error) {
	var config Config

	err := hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
