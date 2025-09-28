// package handlers sets up the handlers for the Web UI
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/alighieri/config"
	"github.com/johannes-kuhfuss/alighieri/dto"
	"github.com/johannes-kuhfuss/alighieri/repositories"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type StatsUiHandler struct {
	Cfg  *config.AppConfig
	Repo *repositories.DefaultDeviceRepository
}

// NewStatsUiHandler creates a new web UI handler and injects its dependencies
func NewStatsUiHandler(cfg *config.AppConfig, repo *repositories.DefaultDeviceRepository) StatsUiHandler {
	return StatsUiHandler{
		Cfg:  cfg,
		Repo: repo,
	}
}

// StatusPage is the handler for the status page
func (uh *StatsUiHandler) StatusPage(c *gin.Context) {
	configData := dto.GetConfig(uh.Cfg)
	c.HTML(http.StatusOK, "status.page.tmpl", gin.H{
		"title":      "Status",
		"configdata": configData,
	})
}

// FileListPage is the handler for the device list page
func (uh *StatsUiHandler) DeviceListPage(c *gin.Context) {
	devices := dto.GetDevices(uh.Repo)
	c.HTML(http.StatusOK, "devicelist.page.tmpl", gin.H{
		"title":   "Device List",
		"devices": devices,
	})
}

// LogsPage is the handler for the page displaying log messages
func (uh *StatsUiHandler) LogsPage(c *gin.Context) {
	logs := logger.GetLogList()
	c.HTML(http.StatusOK, "logs.page.tmpl", gin.H{
		"title": "Logs",
		"logs":  logs,
	})
}

// AboutPage is the handler for the page displaying a short description of the program and its license
func (uh *StatsUiHandler) AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.page.tmpl", gin.H{
		"title": "About",
		"data":  nil,
	})
}
