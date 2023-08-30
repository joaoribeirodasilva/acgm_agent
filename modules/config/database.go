package config

type ConfigDatabase struct {
	Host      string `json:"host" yaml:"host"`
	Port      int    `json:"port" yaml:"port"`
	User      string `json:"user" yaml:"user"`
	Database  string `json:"database" yaml:"database"`
	Password  string `json:"password" yaml:"password"`
	Charset   string `json:"charset" yaml:"charset"`
	ParseTime bool   `json:"parse_time" yaml:"parse_time"`
}
