package config

type ConfigCollect struct {
	Active    bool  `json:"active" yaml:"active"`
	Interval  int64 `json:"collect_interval" yaml:"collect_interval"`
	Aggregate int   `json:"aggregate" yaml:"aggregate"`
}

type CollectProcesses struct {
	ConfigCollect
	ProcessNames []string `json:"process_names" yaml:"process_names"`
}

type ConfigNginx struct {
	Active  bool `json:"active" yaml:"active"`
	Writing bool `json:"writing" yaml:"writing"`
	Waiting bool `json:"waiting" yaml:"waiting"`
}

type HTTPRequests struct {
	Access   bool        `json:"access" yaml:"access"`
	Errors   bool        `json:"errors" yaml:"errors"`
	Duration bool        `json:"duration" yaml:"duration"`
	Status   bool        `json:"status" yaml:"status"`
	Nginx    ConfigNginx `json:"nginx" yaml:"nginx"`
}

type ConfigMetrics struct {
	Host       ConfigCollect    `json:"host" yaml:"host"`
	Cpu        ConfigCollect    `json:"cpu" yaml:"cpu"`
	Memory     ConfigCollect    `json:"memory" yaml:"memory"`
	Disk       ConfigCollect    `json:"disk" yaml:"disk"`
	Network    ConfigCollect    `json:"network" yaml:"network"`
	Load       ConfigCollect    `json:"load" yaml:"load"`
	Partitions ConfigCollect    `json:"partitions" yaml:"partitions"`
	Processes  CollectProcesses `json:"processes" yaml:"processes"`
	HTTP       HTTPRequests     `json:"http" yaml:"http"`
}
