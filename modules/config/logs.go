package config

type ConfigLogs struct {
	Console bool   `json:"console" yaml:"console"`
	Json    bool   `json:"json" yaml:"json"`
	File    bool   `json:"file" yaml:"file"`
	Backups int    `json:"backups" yaml:"backups"`
	Dir     string `json:"dir" yaml:"dir"`
	Name    string `json:"name" yaml:"name"`
	MaxSize int    `json:"max_size" yaml:"max_size"`
	MaxDays int    `json:"max_days" yaml:"max_days"`
}
