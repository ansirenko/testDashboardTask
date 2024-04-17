package testDashboardTask

import (
	"fmt"
	"time"
)

func smartWaiter(action func() error, timeout time.Duration) error {
	timeoutChan := time.After(timeout)
	tick := time.Tick(time.Second)
	var err error
	for {
		select {
		case <-timeoutChan:
			if err != nil {
				return err
			}
			return fmt.Errorf("timeout reached, action failed")
		case <-tick:
			if err = action(); err != nil {
				continue
			}
			return nil
		}
	}
}
