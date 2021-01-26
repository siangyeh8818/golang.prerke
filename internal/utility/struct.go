package utility

import "time"

const (
	//command = `bc -l`
	timeout = 10 * time.Minute
)

type Configs struct {
	User     string          `yaml:"user"`
	Password string          `yaml:"password"`
	Nodes    []AddressConfig `yaml:"nodes"`
	RkeUrl   string          `yaml:"rke_tool"`
}

type AddressConfig struct {
	Address       string
	Info          []string
	Dockerversion string
}

type RkeUserConfig struct {
	Address string
	User    string
}

type RkeConfig struct {
	Nodes []RkeUserConfig
}
