// package app ties together all bits and pieces to start the program
package app

import (
	"time"
)

// initMetrics sets up the Prometheus metrics
func initMetrics() {
}

func updateMetrics() {
	for {
		doUpdate()
		time.Sleep(3 * time.Second)
	}
}

func doUpdate() {
	cfg.RunTime.Mu.Lock()
	defer cfg.RunTime.Mu.Unlock()
}
