package config

type ConfigNginx struct {
	Active  bool `json:"active" yaml:"active"`
	Writing bool `json:"writing" yaml:"writing"`
	Waiting bool `json:"waiting" yaml:"waiting"`
}

type ConfigRequests struct {
	Access   bool        `json:"access" yaml:"access"`
	Errors   bool        `json:"errors" yaml:"errors"`
	Duration bool        `json:"duration" yaml:"duration"`
	Status   bool        `json:"status" yaml:"status"`
	Nginx    ConfigNginx `json:"nginx" yaml:"nginx"`
}

type ConfigMetrics struct {
	CollectInterval int64          `json:"collect_interval" yaml:"collect_interval"`
	Cpu             bool           `json:"cpu" yaml:"cpu"`
	Mem             bool           `json:"mem" yaml:"mem"`
	Disk            bool           `json:"disk" yaml:"disk"`
	Net             bool           `json:"net" yaml:"net"`
	Requests        ConfigRequests `json:"requests" yaml:"requests"`
}
