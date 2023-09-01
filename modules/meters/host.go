package meters

type Host struct {
	Name string `json:"name" yaml:"name"`
	CPUs *CPUs  `json:"cpus" yaml:"cpus"`
}

type Hosts struct {
	Physical   int `json:"physical" yaml:"physical"`
	Containers int `json:"containers" yaml:"containers"`
	Host       []Host
}
