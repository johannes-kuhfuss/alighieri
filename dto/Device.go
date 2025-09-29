// package dto defines the data structures used to exchange information
package dto

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/johannes-kuhfuss/alighieri/domain"
	"github.com/johannes-kuhfuss/alighieri/repositories"
)

// DeviceResp defines the data to be displayed in the device list
type DeviceResp struct {
	Name         string
	FullName     string
	HostName     string
	IPv4         string
	Port         string
	Manufacturer string
	Model        string
	Info         string
	FirstSeen    string
	LastSeen     string
}

// GetDevices retrives all devices maintained in the repository and formats them for display purposes
func GetDevices(repo *repositories.DefaultDeviceRepository) (deviceDta []DeviceResp) {
	if devices := repo.GetAll(); devices != nil {
		for _, device := range *devices {
			dta := DeviceResp{
				Name:         device.Name,
				FullName:     device.FullName,
				HostName:     device.HostName,
				IPv4:         device.IPv4.String(),
				Port:         strconv.Itoa(device.Port),
				Manufacturer: device.Manufacturer,
				Model:        device.Model,
				Info:         combineInfo(device),
				FirstSeen:    device.FirstSeen.Format("2006-01-02 15:04:05"),
				LastSeen:     device.LastSeen.Format("2006-01-02 15:04:05"),
			}
			deviceDta = append(deviceDta, dta)
		}
	}
	sort.SliceStable(deviceDta, func(i, j int) bool {
		if strings.Compare(deviceDta[i].Name, deviceDta[j].Name) > 0 {
			return false
		} else {
			return true
		}
	})
	return
}

func combineInfo(device domain.DeviceInfo) string {
	return fmt.Sprintf("Id: %s, Process: %s, CMCP Version: %s, CMCP Min: %s, Server Version: %s, Channels: %s", device.Id, device.Process, device.CmcpVersion, device.CmcpMin, device.ServerVersion, device.Channels)
}
