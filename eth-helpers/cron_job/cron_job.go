package cron_job

import (
	"time"
)

func GetInterval(interval int) time.Duration {
	if interval == 0 {
		return 100 * time.Second // Default interval
	}
	return time.Duration(interval) * time.Second
}
