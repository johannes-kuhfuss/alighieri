// package service implements the services and their business logic that provide the main part of the program
package service

import (
	"net"
	"strings"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/johannes-kuhfuss/alighieri/config"
	"github.com/johannes-kuhfuss/alighieri/domain"
	"github.com/johannes-kuhfuss/alighieri/repositories"
	"github.com/johannes-kuhfuss/services_utils/logger"
	defaultroute "github.com/nixigaj/go-default-route"
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

func selectNetworkInterface(cfg *config.AppConfig) {
	logger.Info("Determining network interface...")
	if cfg.DeviceScan.InterfaceName != "" {
		logger.Infof("Trying to find interface with name %v", cfg.DeviceScan.InterfaceName)
		iface, err := net.InterfaceByName(cfg.DeviceScan.InterfaceName)
		if err != nil {
			logger.Errorf("Could not find interface with name %v. Using default interface.", cfg.DeviceScan.InterfaceName)
		} else {
			logger.Infof("Found interface with name %v", cfg.DeviceScan.InterfaceName)
			cfg.RunTime.DeviceScanInterface = iface
			return
		}
	}
	defIface, err := defaultroute.DefaultRouteInterface()
	if err != nil {
		logger.Error("Could not find default interface. Giving up...", err)
	} else {
		cfg.RunTime.DeviceScanInterface = defIface
	}
	logger.Infof("Using network interface %v", cfg.RunTime.DeviceScanInterface.Name)
}

// NewDeviceScanService creates a new device scan service and injects its dependencies
func NewDeviceScanService(cfg *config.AppConfig, repo *repositories.DefaultDeviceRepository) DefaultDeviceScanService {
	selectNetworkInterface(cfg)
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
	logger.Infof("Starting device scan run #%v on network interface %v.", s.Cfg.RunTime.DeviceScanNumber, s.Cfg.RunTime.DeviceScanInterface.Name)
	start := time.Now().UTC()
	deviceCount, err := s.scanDevices()
	if err != nil {
		logger.Errorf("Error while scanning for audio devices: %v", err)
	}
	deviceListCount := s.Repo.Size()
	end := time.Now().UTC()
	dur := end.Sub(start)
	logger.Infof("Finished device scan run #%v. Found %v devices. %v device(s) in list total. (%v)", s.Cfg.RunTime.DeviceScanNumber, deviceCount, deviceListCount, dur.String())
	s.Cfg.RunTime.Mu.Lock()
	defer s.Cfg.RunTime.Mu.Unlock()
	s.Cfg.RunTime.DevicesInList = deviceListCount
	s.Cfg.RunTime.DeviceScanRunning = false
	return nil
}

func (s DefaultDeviceScanService) scanDevices() (deviceCount int, err error) {
	var (
		numEntries int
	)
	entriesCh := make(chan *mdns.ServiceEntry, 32)
	go func() {
		for entry := range entriesCh {
			numEntries++
			logger.Infof("Found device %v\r\n", entry.Name)
			device, err := convertEntry(*entry)
			if err != nil {
				logger.Error("Could not convert entry to device", err)
			} else {
				s.storeDevice(device)
			}
		}
	}()
	queryParams := &mdns.QueryParam{
		Service:             s.Cfg.DeviceScan.ServiceName,
		Domain:              "local",
		Timeout:             time.Duration(s.Cfg.DeviceScan.ScanTimeOutSec) * time.Second,
		Interface:           s.Cfg.RunTime.DeviceScanInterface,
		Entries:             entriesCh,
		WantUnicastResponse: false,
	}
	err = mdns.Query(queryParams)
	if err != nil {
		logger.Errorf("Error while querying audio devices: %v", err)
		return 0, err
	}
	close(entriesCh)
	return numEntries, nil
}

func convertEntry(e mdns.ServiceEntry) (dev domain.DeviceInfo, err error) {
	d := domain.DeviceInfo{
		IPv4:      e.AddrV4,
		Port:      e.Port,
		FirstSeen: time.Now(),
		LastSeen:  time.Now(),
	}
	d.FullName = strings.TrimSuffix(e.Name, ".")
	d.HostName = strings.TrimSuffix(e.Host, ".")
	d.Name = shorten(e.Host)
	for _, info := range e.InfoFields {
		if info != "" {
			kvp := strings.Split(info, "=")
			if len(kvp) != 2 {
				logger.Warnf("Could not split %s", info)
				break
			}
			switch strings.ToLower(kvp[0]) {
			case "id":
				d.Id = kvp[1]
			case "process":
				d.Process = kvp[1]
			case "cmcp_vers":
				d.CmcpVersion = kvp[1]
			case "cmcp_min":
				d.CmcpMin = kvp[1]
			case "server_vers":
				d.ServerVersion = kvp[1]
			case "channels":
				d.Channels = kvp[1]
			case "mf":
				d.Manufacturer = kvp[1]
			case "model":
				d.Model = kvp[1]
			}
		}
	}
	return d, nil
}

func shorten(fqdn string) string {
	i := strings.Index(fqdn, ".")
	if i == -1 {
		return fqdn
	}
	return fqdn[0:i]
}

func (s DefaultDeviceScanService) storeDevice(dev domain.DeviceInfo) (err error) {
	oldDev := s.Repo.GetByName(dev.Name)
	if oldDev != nil {
		dev.FirstSeen = oldDev.FirstSeen
	}
	err = s.Repo.Store(dev)
	return err
}
