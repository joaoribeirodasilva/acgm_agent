package base

import (
	"fmt"
	"plugin"
)

type PluginStatus int

const (
	STATUS_UNKNOWN PluginStatus = iota
	STATUS_PLUGIN_NOT_FOUND
	STATUS_DISABLED
	STATUS_STOPPED
	STATUS_STOPPING
	STATUS_STARTING
	STATUS_COLLECTING
	STATUS_SAVING
	STATUS_RUNNING
)

var plugin_status_strings = map[PluginStatus]string{
	STATUS_UNKNOWN:          "unknown",
	STATUS_PLUGIN_NOT_FOUND: "not found",
	STATUS_DISABLED:         "disabled",
	STATUS_STOPPED:          "stopped",
	STATUS_STOPPING:         "stopping",
	STATUS_STARTING:         "starting",
	STATUS_COLLECTING:       "collecting",
	STATUS_SAVING:           "saving",
	STATUS_RUNNING:          "running",
}

func (s PluginStatus) String() string {
	val, ok := plugin_status_strings[s]
	if !ok {
		return plugin_status_strings[STATUS_UNKNOWN]
	}
	return val
}

type Configuration struct {
	host_id   int64
	active    bool
	interval  int64
	aggregate int64
	other     map[string]string
}

type InterfacePlugin interface {
	Factory() error
	Start() error
	Stop() error
	GetData(meter string) (interface{}, error)
	GetStatus() PluginStatus
	SetStatus(status PluginStatus)
	Polling() error
}

type PluginBase struct {
	conf   *Configuration
	plugin InterfacePlugin
}

func (p *PluginBase) Factory(name string, conf *Configuration) error {

	var err error
	p.conf = conf
	err = p.load_plugin(name)
	return err
}

func (p *PluginBase) Start() error {

	if !p.is_loaded() {
		return fmt.Errorf("plugin not found")
	}
	p.plugin.Start()
	return nil
}

func (p *PluginBase) Stop() error {

	if !p.is_loaded() {
		return fmt.Errorf("plugin not found")
	}
	p.plugin.Stop()
	return nil
}

func (p *PluginBase) Polling() (bool, error) {

	var err error
	if !p.is_loaded() {
		return false, fmt.Errorf("plugin not found")
	}

	go func() {
		err = p.plugin.Polling() //poling by active metric
		// store errors (as warning. best effort theory)
		// send an event
	}()

	return false, err
}

func (p *PluginBase) GetData(meter string) (interface{}, error) {
	if !p.is_loaded() {
		return nil, fmt.Errorf("plugin not found")
	}
	return p.plugin.GetData(meter)
}

func (p *PluginBase) GetStatus() PluginStatus {
	if !p.is_loaded() {
		return STATUS_PLUGIN_NOT_FOUND
	}
	return p.plugin.GetStatus()
}

func (p *PluginBase) load_plugin(name string) error {

	mod := fmt.Sprintf("./plugins/%s.so", name)
	plug, err := plugin.Open(mod)
	if err != nil {
		p.plugin = nil
		return err
	}
	temp, err := plug.Lookup("M")
	if err != nil {
		plug = nil
		return err
	}
	var ok bool
	p.plugin, ok = temp.(InterfacePlugin)
	if !ok {
		return fmt.Errorf("plugin %s does not implement InterfacePlugin interface", name)
	}

	return nil
}

func (p *PluginBase) is_loaded() bool {

	return p.plugin != nil
}
