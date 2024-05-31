package utils

import (
	"time"

	"github.com/go-co-op/gocron"
)

var (
	s = gocron.NewScheduler(time.UTC)
)

func CronStart(f func()) {
	return
	s.Every("1m").Do(f)
	s.StartAsync()
}
