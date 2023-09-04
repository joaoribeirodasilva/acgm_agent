package main

import (
	"time"

	"biqx.com.br/acgm_agent/modules/cmd"
	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/database"
	"biqx.com.br/acgm_agent/modules/logger"
	"biqx.com.br/acgm_agent/modules/meters"
)

func main() {

	options, err := cmd.Parse()
	if err != nil {
		panic("ERROR: " + err.Error())
	}

	config := &config.Config{}
	err = config.Read(options)
	if err != nil {
		panic("ERROR: " + err.Error())
	}

	logger.Init(config)
	logger.Log.Info().Msg("Starting ACGM Agent")

	database := database.New(config)
	err = database.Connect()
	if err != nil {
		return
	}

	cpus := meters.NewMetricsCPU(config)
	partitions := meters.NewMetricsPartition(config)
	memory := meters.NewMetricsMemory(config)
	net := meters.NewMetricsNet(config)
	load := meters.NewMetricsLoad(config)
	threadCPUs := meters.NewThread(cpus, config)
	threadPartitions := meters.NewThread(partitions, config)
	threadMemory := meters.NewThread(memory, config)
	threadNet := meters.NewThread(net, config)
	threadLoad := meters.NewThread(load, config)

	threadCPUs.Start()
	threadPartitions.Start()
	threadMemory.Start()
	threadNet.Start()
	threadLoad.Start()

	time.Sleep(5000 * time.Millisecond)

	threadCPUs.Stop()
	threadPartitions.Stop()
	threadMemory.Stop()
	threadNet.Stop()
	threadLoad.Stop()
	database.Disconnect()
	logger.Log.Info().Msg("Terminating ACGM Agent")
}
