package meters

type Meters struct {
	Host      Host      `json:"host" yaml:"host"`
	Processes Processes `json:"processes" yaml:"processes"`
}
