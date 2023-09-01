package meters

import (
	"biqx.com.br/acgm_agent/modules/config"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/cpu"
)

type CPU struct {
	Info    cpu.InfoStat  `json:"info" yaml:"info"`
	Percent float64       `json:"percent" yaml:"percent"`
	Times   cpu.TimesStat `json:"times" yaml:"times"`
	Threads int           `json:"threads" yaml:"threads"`
}

type CPUs struct {
	ThreadName string        `json:"physical" yaml:"physical"`
	Physical   int           `json:"physical" yaml:"physical"`
	Logical    int           `json:"logical" yaml:"logical"`
	Percent    float64       `json:"percent" yaml:"percent"`
	Times      cpu.TimesStat `json:"times" yaml:"times"`
	CPUs       []CPU         `json:"cpus" yaml:"cpus"`
	config     *config.Config
}

func NewCPUs(conf *config.Config) *CPUs {
	c := &CPUs{
		ThreadName: "CPUs",
		config:     conf,
	}
	return c
}

func (c *CPUs) Init() error {

	c.GetCpusCount()
	c.GetCpusInfo()
	return nil
}

func (c *CPUs) GetThreadName() string {
	return c.ThreadName
}

func (c *CPUs) Collect() {
	c.GetCpusTimesTotal()
	c.GetCpusPercentTotal()
	c.GetCpusTimes()
	c.GetCpusPercent()
}

func (c *CPUs) GetCpusCount() {

	physical, err := cpu.Counts(false)
	if err != nil {
		log.Error().Str("namespace", "meters::cpus::GetCpusCount").Err(err).Msg("Failed to get host physical CPU count")
	}
	c.Physical = physical

	logical, err := cpu.Counts(true)
	if err != nil {
		log.Error().Str("namespace", "meters::cpus::GetCpusCount").Err(err).Msg("Failed to get host logical CPU count")
	}
	c.Logical = logical
}

func (c *CPUs) GetCpusInfo() {

	info, err := cpu.Info()
	if err != nil {
		log.Error().Str("namespace", "meters::cpus::GetCpusInfo").Err(err).Msg("Failed to get host CPUs information")
	}

	c.CPUs = []CPU{}
	for idx, item := range info {
		ncpu := CPU{
			Info:    item,
			Threads: c.Logical / c.Physical,
		}
		c.CPUs = append(c.CPUs, ncpu)
		log.Debug().Str("namespace", "meters::cpus::GetCpusInfo").Msgf("Found CPU[%d] %s", idx, item.String())
	}

}

func (c *CPUs) GetCpusTimesTotal() {

	total, err := cpu.Times(false)
	if err != nil {
		log.Error().Str("namespace", "meters::cpus::GetCpusTimesTotal").Err(err).Msg("Failed to get host CPUs total times")
	}

	c.Times = total[0]
}

func (c *CPUs) GetCpusTimes() {

	log.Debug().Str("namespace", "meters::cpus::GetCpusTimes").Msg("Enter")
	times, err := cpu.Times(true)
	if err != nil {
		log.Error().Str("namespace", "meters::cpus::GetCpusTimes").Err(err).Msg("Failed to get host CPUs times")
	}

	for idx, item := range times {
		c.CPUs[idx].Times = item
	}
}

func (c *CPUs) GetCpusPercentTotal() {

	percent, err := cpu.Percent(0, false)
	if err != nil {
		log.Error().Str("namespace", "meters::cpus::GetCpusPercentTotal").Err(err).Msg("Failed to get host CPUs percent total")
	}

	c.Percent = percent[0]
}

func (c *CPUs) GetCpusPercent() {

	percent, err := cpu.Percent(0, true)
	if err != nil {
		log.Error().Str("namespace", "meters::cpus::GetCpusPercent").Err(err).Msg("Failed to get host CPUs percent")
	}

	for idx, item := range percent {
		c.CPUs[idx].Percent = item
		log.Debug().Str("namespace", "meters::cpus::GetCpusPercent").Msgf("CPU[%d] Percent: %f%%", idx, item)
	}
}
