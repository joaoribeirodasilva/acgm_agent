package config

type CollectInterval struct {
	CollectInterval int64 `json:"collect_interval" yaml:"collect_interval"`
	Aggregate       int   `json:"aggregate" yaml:"aggregate"`
}

type CollectProcesses struct {
	CollectInterval int64    `json:"collect_interval" yaml:"collect_interval"`
	Aggregate       int      `json:"aggregate" yaml:"aggregate"`
	ProcessNames    []string `json:"process_names" yaml:"process_names"`
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
	Cpu        CollectInterval  `json:"cpu" yaml:"cpu"`
	Memory     CollectInterval  `json:"memory" yaml:"memory"`
	Disk       CollectInterval  `json:"disk" yaml:"disk"`
	Network    CollectInterval  `json:"network" yaml:"network"`
	Load       CollectInterval  `json:"load" yaml:"load"`
	Partitions CollectInterval  `json:"partitions" yaml:"partitions"`
	Processes  CollectProcesses `json:"processes" yaml:"processes"`
	HTTP       HTTPRequests     `json:"http" yaml:"http"`
}
