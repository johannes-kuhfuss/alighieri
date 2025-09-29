// package service implements the services and their business logic that provide the main part of the program
package service

import (
	"time"

	"github.com/johannes-kuhfuss/alighieri/config"
	"github.com/johannes-kuhfuss/alighieri/repositories"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type DeviceScanService interface {
	Scan()
	ScanRun() error
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

func (s DefaultDeviceScanService) Scan() {
	for s.Cfg.DeviceScan.DeviceScanRun {
		s.ScanRun()
		time.Sleep(time.Duration(s.Cfg.DeviceScan.ScanCycleSec) * time.Second)
	}
}

// Scan orchestrates the process of querying audio devices and adding the retrieved information to the device repository
func (s DefaultDeviceScanService) ScanRun() error {
	s.Cfg.RunTime.DeviceScanNumber++
	s.Cfg.RunTime.LastDeviceScanDate = time.Now()
	s.Cfg.RunTime.DeviceScanRunning = true
	logger.Infof("Starting device scan run #%v on network interface XXX.", s.Cfg.RunTime.DeviceScanNumber)
	start := time.Now().UTC()
	deviceCount, err := s.scanDevices()
	if err != nil {
		logger.Errorf("Error while scanning for audio devices: %v", err)
	}
	deviceListCount := s.Repo.Size()
	end := time.Now().UTC()
	dur := end.Sub(start)
	logger.Infof("Finished device scan run #%v. Found %v devices. %v file(s) in list total. (%v)", s.Cfg.RunTime.DeviceScanNumber, deviceCount, deviceListCount, dur.String())
	s.Cfg.RunTime.Mu.Lock()
	defer s.Cfg.RunTime.Mu.Unlock()
	s.Cfg.RunTime.DevicesInList = deviceListCount
	s.Cfg.RunTime.DeviceScanRunning = false
	return nil
}

func (s DefaultDeviceScanService) scanDevices() (deviceCount int, err error) {
	return 0, nil
}
