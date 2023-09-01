package logger

import (
	"biqx.com.br/acgm_agent/modules/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Initializes the application log system.
// Must be called right after the configuration read.
func Init(conf *config.Config) error {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	level := zerolog.InfoLevel
	strLevel := "info"
	if conf.Settings.Debug {
		level = zerolog.DebugLevel
		strLevel = "debug"
	}
	zerolog.SetGlobalLevel(level)
	log.Info().Msgf("Log service initialized with (%s) level", strLevel)

	// TODO: Add support for log file

	return nil

}
