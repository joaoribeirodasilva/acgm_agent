package meters

import (
	"fmt"
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/logger"
	"github.com/shirou/gopsutil/net"
)

var NET_THREAD_NAME = "network"

type Interface struct {
	Interface net.InterfaceStat  `json:"interface" yaml:"interface"`
	IOCounter net.IOCountersStat `json:"io_counter" yaml:"io_counter"`
}

type Net struct {
	DateTime    time.Time               `json:"date_time" yaml:"date_time"`
	Interfaces  []Interface             `json:"interfaces" yaml:"interfaces"`
	Connections []net.ConnectionStat    `json:"connections" yaml:"connections"`
	Protocols   []net.ProtoCountersStat `json:"protocols" yaml:"protocols"`
}

type MetricsNet struct {
	Metrics     []Net          `json:"metrics" yaml:"metrics"`
	config      *config.Config `json:"-" yaml:"-"`
	init_failed bool           `json:"-" yaml:"-"`
	collecting  bool           `json:"-" yaml:"-"`
	cutting     bool           `json:"-" yaml:"-"`
}

func NewMetricsNet(conf *config.Config) *MetricsNet {
	nm := &MetricsNet{
		config:      conf,
		init_failed: false,
		collecting:  false,
		cutting:     false,
	}
	return nm
}

func (nm *MetricsNet) Init() error {

	nm.init_failed = false

	if err := nm.Collect(); err != nil {
		nm.init_failed = true
		return err
	}

	return nil
}

func (nm *MetricsNet) IsInitFailed() bool {
	return nm.init_failed
}

func (nm *MetricsNet) GetThreadName() string {
	return NET_THREAD_NAME
}

func (nm *MetricsNet) Collect() error {
	var start, end, diff int64
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}
	if nm.config.Settings.Debug {
		start = time.Now().UnixNano() / int64(time.Millisecond)
	}
	nm.collecting = true
	for nm.cutting {
		logger.Log.Debug().Msg("waiting for cut to finish")
		time.Sleep(50 * time.Millisecond)
	}
	logger.Log.Debug().Msg("collecting data")

	metrics := Net{}
	metrics.DateTime = time.Now()

	if err := nm.GetInterfaces(&metrics); err != nil {
		return err
	}

	if err := nm.GetIOCounters(&metrics); err != nil {
		return err
	}

	if err := nm.GetConnections(&metrics); err != nil {
		return err
	}

	if err := nm.GetProtocols(&metrics); err != nil {
		return err
	}

	nm.Metrics = append(nm.Metrics, metrics)

	logger.Log.Debug().Msg("finish data collection for network")
	if nm.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		logger.Log.Debug().Msgf("collect took %d ms", diff)
	}
	nm.collecting = false
	return nil
}

func (nm *MetricsNet) Cut() (*[]Net, error) {

	var start, end, diff int64
	if nm.init_failed {
		return nil, nil
	}
	if nm.config.Settings.Debug {
		start = time.Now().UnixNano() / int64(time.Millisecond)
	}
	nm.cutting = true
	for nm.collecting {
		logger.Log.Debug().Msg("waiting for collect to finish")
		time.Sleep(50 * time.Millisecond)
	}
	logger.Log.Debug().Msg("cutting data")

	metrics := nm.Metrics
	nm.Metrics = []Net{}
	logger.Log.Debug().Msgf("finish data cut with %d metrics", len(metrics))

	if nm.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		logger.Log.Debug().Msgf("cut took %d ms", diff)
	}
	nm.cutting = false

	return &metrics, nil
}

func (nm *MetricsNet) GetInterfaces(n *Net) error {

	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}

	ifaces, err := net.Interfaces()
	if err != nil || len(ifaces) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host interfaces data")
		return err
	}

	for _, stats := range ifaces {
		iface := Interface{Interface: stats}
		n.Interfaces = append(n.Interfaces, iface)
		// fmt.Printf("Interface: %+v\n", iface)
	}

	return nil
}

func (nm *MetricsNet) GetIOCounters(n *Net) error {

	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}

	counters, err := net.IOCounters(true)
	if err != nil || len(counters) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host interfaces IO counters data")
		return err
	}

	for _, counter := range counters {
		for idx, iface := range n.Interfaces {
			if iface.Interface.Name == counter.Name {
				n.Interfaces[idx].IOCounter = counter
			}
		}
	}

	// fmt.Printf("Interface: %+v\n", n.Interfaces)

	return nil
}

func (nm *MetricsNet) GetConnections(n *Net) error {

	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}

	connections, err := net.Connections("all")
	if err != nil || len(connections) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host interfaces connections data")
		return err
	}

	for _, connection := range connections {
		if connection.Status == "NONE" {
			continue
		}

		n.Connections = append(n.Connections, connection)

		// fmt.Printf("Connection: %+v\n", connection)
	}

	return nil
}

func (nm *MetricsNet) GetProtocols(n *Net) error {

	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}

	protocols, err := net.ProtoCounters([]string{"ip", "icmp", "icmpmsg", "tcp", "udp", "udplite"})
	if err != nil || len(protocols) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host protocols data")
		return err
	}

	n.Protocols = protocols

	return nil
}

func (nm *MetricsNet) error_init() error {
	str_err := "initialization failed, init must be run again"
	err := fmt.Errorf(str_err)
	return err
}
