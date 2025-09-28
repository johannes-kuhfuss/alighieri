// package dto defines the data structures used to exchange information
package dto

import (
	"strconv"
	"time"

	"github.com/johannes-kuhfuss/alighieri/config"
	"github.com/robfig/cron/v3"
)

// ConfigResp converted configuration data for display on the web UI
type ConfigResp struct {
	ServerHost                 string
	ServerPort                 string
	ServerTlsPort              string
	ServerGracefulShutdownTime string
	ServerUseTls               string
	ServerCertFile             string
	ServerKeyFile              string
	GinMode                    string
	StartDate                  string
	LogFile                    string
}

// setStartDate sets the service start date and adds the run duration
func setStartDate(date time.Time) string {
	dur := time.Since(date)
	return convertDate(date) + " (running for " + dur.String() + ")"
}

// convertDate converts a date to its display format
func convertDate(date time.Time) string {
	if date.IsZero() {
		return "N/A"
	} else {
		return date.Local().Format("2006-01-02 15:04:05 -0700 MST")
	}
}

// getNextJobDate retrieves a job's next execution date and returns it in its display format
func getNextJobDate(cfg *config.AppConfig, jobId int) string {
	if cfg.RunTime.BgJobs != nil && cfg.RunTime.BgJobs.Entry(cron.EntryID(jobId)).Valid() {
		return cfg.RunTime.BgJobs.Entry(cron.EntryID(jobId)).Next.String()
	} else {
		return "N/A"
	}
}

func formatLogFile(logFile string) string {
	if logFile == "" {
		return "Logging to file disabled"
	}
	return logFile
}

// GetConfig converts the configuration to its display format
func GetConfig(cfg *config.AppConfig) (resp ConfigResp) {
	cfg.RunTime.Mu.Lock()
	defer cfg.RunTime.Mu.Unlock()
	resp = ConfigResp{
		ServerHost:                 cfg.Server.Host,
		ServerPort:                 cfg.Server.Port,
		ServerTlsPort:              cfg.Server.TlsPort,
		ServerGracefulShutdownTime: strconv.Itoa(cfg.Server.GracefulShutdownTime),
		ServerUseTls:               strconv.FormatBool(cfg.Server.UseTls),
		ServerCertFile:             cfg.Server.CertFile,
		ServerKeyFile:              cfg.Server.KeyFile,
		GinMode:                    cfg.Gin.Mode,
		LogFile:                    formatLogFile(cfg.Server.LogFile),
	}
	resp.StartDate = setStartDate(cfg.RunTime.StartDate)
	if cfg.Server.Host == "" {
		resp.ServerHost = "localhost"
	}
	return
}
