package config

import (
	"github.com/alash3al/xyr/driver"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	DataDir string   `hcl:"data_dir"`
	Tables  []*Table `hcl:"table,block"`
}

type Table struct {
	Name    string   `hcl:"name,label"`
	DSN     string   `hcl:"dsn"`
	Loader  string   `hcl:"loader"`
	Columns []string `hcl:"columns"`

	DriverInstance driver.Driver
}

func LoadConfigFromFile(filename string) (*Config, error) {
	var config Config

	if err := hclsimple.DecodeFile(filename, nil, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
