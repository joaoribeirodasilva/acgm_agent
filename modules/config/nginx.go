package config

type Nginx struct {
	SitesAvailable string `json:"sites_available" yaml:"sites_available"`
	SitesEnabled   string `json:"sites_enabled" yaml:"sites_enabled"`
	Reload         string `json:"reload" yaml:"reload"`
	Test           string `json:"test" yaml:"test"`
	Logs           string `json:"logs" yaml:"logs"`
	Error          string `json:"error" yaml:"error"`
	Access         string `json:"access" yaml:"access"`
	DefaultRoot    string `json:"default_root" yaml:"default_root"`
	DefaultFile    string `json:"default_file" yaml:"default_file"`
	SslDir         string `json:"ssl_dir" yaml:"ssl_dir"`
}
