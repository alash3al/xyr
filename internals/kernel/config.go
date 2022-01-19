package kernel

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	DataDir      string   `hcl:"data_dir"`
	Tables       []*Table `hcl:"table,block"`
	Debug        bool     `hcl:"debug"`
	WorkersCount int      `hcl:"workers_count"`
}

func LoadConfigFromFile(filename string) (*Config, error) {
	var config Config

	if err := hclsimple.DecodeFile(filename, nil, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
