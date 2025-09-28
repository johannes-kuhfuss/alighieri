package repositories

import (
	"testing"

	"github.com/johannes-kuhfuss/alighieri/config"
	"github.com/johannes-kuhfuss/alighieri/domain"
	"github.com/stretchr/testify/assert"
)

var (
	cfg  config.AppConfig
	repo DefaultDeviceRepository
)

func setupTest() {
	repo = NewDeviceRepository(&cfg)
}

func TestNewDeviceRepositoryCreatesEmptyList(t *testing.T) {
	setupTest()
	assert.EqualValues(t, 0, repo.Size())
}

func TestGetByNameEmptyListReturnsNil(t *testing.T) {
	setupTest()
	res := repo.GetByName("B")
	assert.Nil(t, res)
}

func TestGetAllEmptyListReturnsNil(t *testing.T) {
	setupTest()
	res := repo.GetAll()
	assert.Nil(t, res)
}

func TestStoreItemWithEmptyNameReturnsError(t *testing.T) {
	setupTest()
	di := domain.DeviceInfo{}
	err := repo.Store(di)
	assert.NotNil(t, err)
	assert.EqualValues(t, "cannot add item with empty name to list", err.Error())
}

func TestGetByPathCorrectPathReturnsElement(t *testing.T) {
	setupTest()
	di := domain.DeviceInfo{
		Name: "A",
	}
	err := repo.Store(di)
	res := repo.GetByName("A")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.EqualValues(t, "A", res.Name)
}

func TestGetAllReturnsAllElements(t *testing.T) {
	setupTest()
	di1 := domain.DeviceInfo{
		Name: "A",
	}
	di2 := domain.DeviceInfo{
		Name: "B",
	}
	repo.Store(di1)
	repo.Store(di2)
	size := repo.Size()
	res := repo.GetAll()
	assert.NotNil(t, size)
	assert.EqualValues(t, 2, size)
	assert.EqualValues(t, 2, len(*res))
}

func TestDeleteNonExistingElementReturnsError(t *testing.T) {
	setupTest()
	err := repo.Delete("A")
	assert.NotNil(t, err)
	assert.EqualValues(t, "item with name A does not exist", err.Error())
}

func TestDeleteExistingElementDeletesElement(t *testing.T) {
	setupTest()
	di := domain.DeviceInfo{
		Name: "A",
	}
	repo.Store(di)
	sizeBefore := repo.Size()
	err := repo.Delete("A")
	sizeAfter := repo.Size()
	assert.Nil(t, err)
	assert.EqualValues(t, 1, sizeBefore)
	assert.EqualValues(t, 0, sizeAfter)
}

func TestDeleteAllDeletesAllElements(t *testing.T) {
	setupTest()
	di1 := domain.DeviceInfo{
		Name: "A",
	}
	di2 := domain.DeviceInfo{
		Name: "B",
	}
	repo.Store(di1)
	repo.Store(di2)
	sizeBefore := repo.Size()
	repo.DeleteAllData()
	sizeAfter := repo.Size()
	assert.EqualValues(t, 2, sizeBefore)
	assert.EqualValues(t, 0, sizeAfter)
}
