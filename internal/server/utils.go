package server

import (
	"time"
)

func retry(attempts int, interval time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(interval)
	}
	return err
}
