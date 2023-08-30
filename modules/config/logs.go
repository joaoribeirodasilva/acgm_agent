package config

type ConfigLogs struct {
	Retention string `json:"retention" yaml:"retention"`
	Dir       string `json:"dir" yaml:"dir"`
	File      string `json:"file" yaml:"file"`
	MaxSize   string `json:"max_size" yaml:"max_size"`
}
