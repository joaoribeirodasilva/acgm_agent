package config

type ConfigSettings struct {
	Service  bool `json:"service" yaml:"service"`
	Ipv4Only bool `json:"ipv4only" yaml:"ipv4only"`
	Debug    bool `json:"debug" yaml:"debug"`
}
