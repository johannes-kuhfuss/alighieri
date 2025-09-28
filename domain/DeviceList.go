// package domain defines the core data structures
package domain

import (
	"sync"
)

// DeviceInfo defines the information maintained per device entry
type DeviceInfo struct {
	Name string
}

type DeviceList []DeviceInfo

// SafeDeviceList adds a mutex to allow thread-safe access of the file data entries
type SafeDeviceList struct {
	sync.RWMutex
	Devices map[string]DeviceInfo
}
