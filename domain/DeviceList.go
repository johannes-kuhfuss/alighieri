// package domain defines the core data structures
package domain

import (
	"net"
	"sync"
	"time"
)

// DeviceInfo defines the information maintained per device entry
type DeviceInfo struct {
	Name          string
	FullName      string
	HostName      string
	IPv4          net.IP
	Port          int
	Id            string
	Process       string
	CmcpVersion   string
	CmcpMin       string
	ServerVersion string
	Channels      string
	Manufacturer  string
	Model         string
	FirstSeen     time.Time
	LastSeen      time.Time
}

type DeviceList []DeviceInfo

// SafeDeviceList adds a mutex to allow thread-safe access of the file data entries
type SafeDeviceList struct {
	sync.RWMutex
	Devices map[string]DeviceInfo
}
