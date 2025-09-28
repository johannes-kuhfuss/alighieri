// package dto defines the data structures used to exchange information
package dto

import (
	"sort"
	"strings"

	"github.com/johannes-kuhfuss/alighieri/repositories"
)

// DeviceResp defines the data to be displayed in the file list
type DeviceResp struct {
	Name string
}

/*
func formatTime(t1 time.Time) string {
	if t1.IsZero() {
		return "N/A"
	}
	return t1.Format("15:04")
}
*/

// GetDevices retrives all devices maintained in the repository and formats them for display purposes
func GetDevices(repo *repositories.DefaultDeviceRepository) (deviceDta []DeviceResp) {
	if devices := repo.GetAll(); devices != nil {
		for _, device := range *devices {
			dta := DeviceResp{
				Name: device.Name,
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
