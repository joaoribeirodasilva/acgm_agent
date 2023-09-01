package logger

import (
	"io"
	"os"
	"path"

	"biqx.com.br/acgm_agent/modules/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zerolog.Logger

// Initializes the application log system.
// Must be called right after the configuration read.
func Init(conf *config.Config) error {
	var writers []io.Writer
	if conf.Logs.Console {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if conf.Logs.File {
		writers = append(writers, rollingFile(conf))
	}
	mw := io.MultiWriter(writers...)

	level := zerolog.InfoLevel
	if conf.Settings.Debug {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	l := zerolog.New(mw).With().Timestamp().Caller().Logger()
	Log = &l

	Log.Info().Msgf("Log service initialized with (%d) level", level)

	return nil

}

func rollingFile(conf *config.Config) io.Writer {
	if err := os.MkdirAll(conf.Logs.Dir, 0744); err != nil {
		log.Error().Err(err).Str("path", conf.Logs.Dir).Msg("can't create log directory")
		panic(err)
	}

	return &lumberjack.Logger{
		Filename:   path.Join(conf.Logs.Dir, conf.Logs.Name),
		MaxBackups: conf.Logs.Backups, // files
		MaxSize:    conf.Logs.MaxSize, // megabytes
		MaxAge:     conf.Logs.MaxDays, // days
	}
}
