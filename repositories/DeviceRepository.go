// Package repositories implements an in-memory store for representing the data of the files scanned
package repositories

import (
	"errors"
	"fmt"

	"github.com/johannes-kuhfuss/alighieri/config"
	"github.com/johannes-kuhfuss/alighieri/domain"
)

type DeviceRepository interface {
	Exists(string) bool
	Size() int
	GetByName(string) *domain.DeviceInfo
	GetAll() *domain.DeviceList
	Store(domain.DeviceInfo) error
	Delete(string) error
	DeleteAllData()
}

type DefaultDeviceRepository struct {
	Cfg *config.AppConfig
}

var (
	deviceList domain.SafeDeviceList
)

// NewDeviceRepository creates a new device repository. You need to pass in the configuration
func NewDeviceRepository(cfg *config.AppConfig) DefaultDeviceRepository {
	deviceList.Devices = make(map[string]domain.DeviceInfo)
	return DefaultDeviceRepository{
		Cfg: cfg,
	}
}

// Exists checks whether a device identified by its name exists in the repository
func (dr DefaultDeviceRepository) Exists(name string) bool {
	deviceList.RLock()
	defer deviceList.RUnlock()
	_, ok := deviceList.Devices[name]
	return ok
}

// Size returns the number of devices stored in the repository
func (dr DefaultDeviceRepository) Size() int {
	deviceList.RLock()
	defer deviceList.RUnlock()
	return len(deviceList.Devices)
}

// GetByName returns a device's information where the device is identified by its name. If no device matches, the method returns nil
func (dr DefaultDeviceRepository) GetByName(name string) *domain.DeviceInfo {
	var di domain.DeviceInfo
	if !dr.Exists(name) {
		return nil
	}
	deviceList.RLock()
	defer deviceList.RUnlock()
	di = deviceList.Devices[name]
	return &di
}

// GetAll returns all device data from the repository. Returns nil if repository is empty
func (dr DefaultDeviceRepository) GetAll() *domain.DeviceList {
	var list domain.DeviceList
	if dr.Size() == 0 {
		return nil
	}
	deviceList.RLock()
	defer deviceList.RUnlock()
	for _, device := range deviceList.Devices {
		list = append(list, device)
	}
	return &list
}

// Store stores a new device information entry into the repository
func (dr DefaultDeviceRepository) Store(di domain.DeviceInfo) error {
	if di.Name == "" {
		return errors.New("cannot add item with empty name to list")
	}
	deviceList.Lock()
	defer deviceList.Unlock()
	deviceList.Devices[di.Name] = di
	return nil
}

// Delete a device information entry from the repository, if it exists
func (dr DefaultDeviceRepository) Delete(name string) error {
	if !dr.Exists(name) {
		return fmt.Errorf("item with name %v does not exist", name)
	}
	deviceList.Lock()
	defer deviceList.Unlock()
	delete(deviceList.Devices, name)
	return nil
}

// DeleteAllData removes all entries from the repository
func (dr DefaultDeviceRepository) DeleteAllData() {
	deviceList.Lock()
	defer deviceList.Unlock()
	deviceList.Devices = make(map[string]domain.DeviceInfo)
}
