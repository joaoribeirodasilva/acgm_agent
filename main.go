package main

import (
	"time"

	"biqx.com.br/acgm_agent/modules/cmd"
	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/database"
	"biqx.com.br/acgm_agent/modules/logger"
	"biqx.com.br/acgm_agent/modules/meters"
	"github.com/rs/zerolog/log"
)

func main() {

	options, err := cmd.Parse()
	if err != nil {
		panic("ERROR: " + err.Error())
	}

	// fmt.Printf("%+v, %s, %t\n", options, *options.ConfigFile, *options.Service)

	config := &config.Config{}
	err = config.Read(options)
	if err != nil {
		panic("ERROR: " + err.Error())
	}

	logger.Init(config)
	log.Info().Str("package", "Main").Err(err).Msg("Starting ACGM Agent")

	database := database.New(config)
	err = database.Connect()
	if err != nil {
		return
	}

	cpus := meters.NewCPUs(config)
	cpus.Init()
	thread := meters.NewThread(cpus, config)
	err = thread.Start()
	if err != nil {
		database.Disconnect()
	}
	time.Sleep(5000 * time.Millisecond)
	thread.Stop()

	// meters := meters.NewProcesses()
	// meters.Start(config)

	// time.Sleep(5000 * time.Millisecond)

	// meters.Stop()
	database.Disconnect()
	log.Info().Str("package", "Main").Err(err).Msg("Terminating ACGM Agent")
}
