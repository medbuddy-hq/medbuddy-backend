package utility

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func NewLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339,
		FullTimestamp:    true,
		ForceColors:      true,
	})

	return logger
}
