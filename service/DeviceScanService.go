// package service implements the services and their business logic that provide the main part of the program
package service

import (
	"github.com/johannes-kuhfuss/alighieri/config"
	"github.com/johannes-kuhfuss/alighieri/repositories"
)

type DeviceScanService interface {
	Scan() error
}

// The DeviceScan service scans for available audio devices
type DefaultDeviceScanService struct {
	Cfg  *config.AppConfig
	Repo *repositories.DefaultDeviceRepository
}

// NewDeviceScanService creates a new device scan service and injects its dependencies
func NewDeviceScanService(cfg *config.AppConfig, repo *repositories.DefaultDeviceRepository) DefaultDeviceScanService {
	return DefaultDeviceScanService{
		Cfg:  cfg,
		Repo: repo,
	}
}

// Scan orchestrates the process of querying audio devices and adding the retrieved information to the device repository
func (s DefaultDeviceScanService) Scan() error {
	return nil
}
