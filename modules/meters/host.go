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

type Host struct {
	Name        string         `json:"name" yaml:"name"`
	CPUs        *CPUs          `json:"cpus" yaml:"cpus"`
	Memory      *Memory        `json:"memory" yaml:"memory"`
	Load        *Load          `json:"load" yaml:"load"`
	Net         *Net           `json:"network" yaml:"network"`
	Partitions  *Partitions    `json:"disk" yaml:"disk"`
	Processes   *Processes     `json:"processes" yaml:"processes"`
	config      *config.Config `json:"-" yaml:"-"`
	init_failed bool           `json:"-" yaml:"-"`
	collecting  bool           `json:"-" yaml:"-"`
	cutting     bool           `json:"-" yaml:"-"`
	events      HostEvents     `json:"-" yaml:"-"`
}

type HostEvents struct {
	Listen []string
	Signal map[string]goevents.Event
}

func (h *Host) Init() error {
	// Transfer events to each class signal
	listens := []string{
		"meter.cpus.collected",
		"meter.cpus.aggregated",
		"meter.cpus.inserted",
		"meter.cpus.initialized",
		"meter.memory.collected",
		"meter.memory.aggregated",
		"meter.memory.inserted",
		"meter.load.collected",
		"meter.load.aggregated",
		"meter.load.inserted",
		"meter.network.collected",
		"meter.network.aggregated",
		"meter.network.inserted",
		"meter.disk.collected",
		"meter.disk.aggregated",
		"meter.disk.inserted",
		"meter.processes.collected",
		"meter.processes.aggregated",
		"meter.processes.inserted",
	}

	signals := []string{
		"meter.host.start",
		"meter.host.stop",
		"meter.host.initialized",
	}

	h.events.Signal = make(map[string]goevents.Event)
	for _, signal := range signals {
		h.events.Signal = map[string]goevents.Event{
			signal: goevents.New(signal),
		}
	}

	return nil
}

func (h *Host) Start() error {
	return nil
}

func (h *Host) Stop() error {
	return nil
}

func (h *Host) Init() error {

	h.init_failed = false
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
