package meters

import (
	"biqx.com.br/acgm_agent/modules/config"
	goevents "github.com/jonhoo/go-events"
)

var HOST_THREAD_NAME = "host"

// Event based metrics
// package main

// import (
// 	"fmt"

// 	goevents
// )

// func main() {

// 	chn := event.Listen("example.hello")
// 	fmt.Println("listening for event")

// 	go func() {
// 		fmt.Println("firing event")
// 		event.Signal("example.hello")
// 		fmt.Println("event fired")
// 	}()

// 	fmt.Println("waiting for channel")
// 	e := <-chn

// 	fmt.Println(e.Tag)
// 	fmt.Println("exiting")
// }

type Signal string

const (
	SIGNAL_HOST_STARTING    Signal = "meter.host.starting"
	SIGNAL_HOST_STARTED     Signal = "meter.host.started"
	SIGNAL_HOST_INITIALIZED Signal = "meter.host.initialized"
	SIGNAL_HOST_STOPPING    Signal = "meter.host.stopping"
	SIGNAL_HOST_STOPPED     Signal = "meter.host.stopped"
)

type Host struct {
	Name            string          `json:"name" yaml:"name"`
	CPUs            *CPUs           `json:"cpus" yaml:"cpus"`
	Memory          *Memory         `json:"memory" yaml:"memory"`
	Load            *Load           `json:"load" yaml:"load"`
	Net             *Net            `json:"network" yaml:"network"`
	Partitions      *Partitions     `json:"disk" yaml:"disk"`
	Processes       *Processes      `json:"processes" yaml:"processes"`
	config          *config.Config  `json:"-" yaml:"-"`
	init_failed     bool            `json:"-" yaml:"-"`
	collecting      bool            `json:"-" yaml:"-"`
	cutting         bool            `json:"-" yaml:"-"`
	received_events map[string]bool `json:"-" yaml:"-"`
}

type HostEvents struct {
	Received map[string]bool
}

func (h *Host) Init() error {
	// Transfer events to each class signal
	go h.ListenMetrics()
	return nil
}

func (h *Host) Start() error {
	return nil
}

func (h *Host) Stop() error {
	return nil
}

func (h *Host) IsInitFailed() bool {
	return h.init_failed
}

func (h *Host) GetThreadName() string {
	return HOST_THREAD_NAME
}

func (h *Host) Collect() error {
	return nil
}

func (h *Host) ListenMetrics() error {
	// Loop until meter.cpu.stopped
	// The events must be statis
	chns := goevents.Listen("metrics.")
	for event := range chns {
		switch event.Tag {
		case "meter.cpu":
			break
		case "meter.memory":
			break
		case "meter.load":
			break
		case "meter.network":
			break
		case "meter.disk":
			break
		case "meter.processes":
			break
		case "meter.nginx":
			break
		}
	}
	return nil
}
