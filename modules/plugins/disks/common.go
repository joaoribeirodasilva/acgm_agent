package disks

type ChangeStatusFunction func(status int) error

type Meter struct {
	name    string
	metrics []string
	conf    Configuration
	model   interface{}
	status  PluginStatus
}
