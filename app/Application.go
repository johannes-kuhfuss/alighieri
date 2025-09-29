// package app ties together all bits and pieces to start the program
package app

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/alighieri/config"
	"github.com/johannes-kuhfuss/alighieri/handlers"
	"github.com/johannes-kuhfuss/alighieri/repositories"
	"github.com/johannes-kuhfuss/alighieri/service"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
)

var (
	cfg            config.AppConfig
	server         http.Server
	appEnd         chan os.Signal
	ctx            context.Context
	cancel         context.CancelFunc
	statsUiHandler handlers.StatsUiHandler
	deviceRepo     repositories.DefaultDeviceRepository
	scanService    service.DefaultDeviceScanService
)

// StartApp orchestrates the startup of the application
func StartApp() {
	getCmdLine()
	err := config.InitConfig(config.EnvFile, &cfg)
	if err != nil {
		panic(err)
	}
	logger.Init(cfg.Server.LogFile)
	logger.Info("Starting application...")
	if cfg.Server.LogFile != "" {
		logger.Infof("Logging to file: %v", cfg.Server.LogFile)
	} else {
		logger.Info("Logging to file disabled")
	}
	initRouter()
	initServer()
	initMetrics()
	wireApp()
	mapUrls()
	RegisterForOsSignals()
	scheduleBgJobs()
	go startServer()
	go updateMetrics()
	go scanService.Scan()

	<-appEnd
	cleanUp()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed", err)
	} else {
		logger.Info("Graceful shutdown finished")
	}
}

// getCmdLine checks the command line arguments
func getCmdLine() {
	flag.StringVar(&config.EnvFile, "config.file", ".env", "Specify location of config file. Default is .env")
	flag.Parse()
}

// initRouter initializes gin-gonic as the router
func initRouter() {
	gin.SetMode(cfg.Gin.Mode)
	router := gin.New()
	if cfg.Gin.LogToLogger {
		gin.DefaultWriter = logger.GetLogger()
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.SetTrustedProxies(nil)
	globPath := filepath.Join(cfg.Gin.TemplatePath, "*.tmpl")
	router.LoadHTMLGlob(globPath)

	cfg.RunTime.Router = router
}

// initServer checks whether https is enabled and initializes the web server accordingly
func initServer() {
	var tlsConfig tls.Config

	if cfg.Server.UseTls {
		tlsConfig = tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS13,
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
				tls.CurveP384,
			},
		}
	}
	if cfg.Server.UseTls {
		cfg.RunTime.ListenAddr = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.TlsPort)
	} else {
		cfg.RunTime.ListenAddr = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	}

	server = http.Server{
		Addr:              cfg.RunTime.ListenAddr,
		Handler:           cfg.RunTime.Router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    0,
	}
	if cfg.Server.UseTls {
		server.TLSConfig = &tlsConfig
		server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	}
}

// wireApp initializes the services in the right order and injects the dependencies
func wireApp() {
	deviceRepo = repositories.NewDeviceRepository(&cfg)
	statsUiHandler = handlers.NewStatsUiHandler(&cfg, &deviceRepo)
	scanService = service.NewDeviceScanService(&cfg, &deviceRepo)
}

// mapUrls defines the handlers for the available URLs
func mapUrls() {
	cfg.RunTime.Router.GET("/", statsUiHandler.StatusPage)
	cfg.RunTime.Router.GET("/devicelist", statsUiHandler.DeviceListPage)
	cfg.RunTime.Router.GET("/logs", statsUiHandler.LogsPage)
	cfg.RunTime.Router.GET("/about", statsUiHandler.AboutPage)
	cfg.RunTime.Router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// RegisterForOsSignals listens for OS signals terminating the program and sends an internal signal to start cleanup
func RegisterForOsSignals() {
	appEnd = make(chan os.Signal, 1)
	signal.Notify(appEnd, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

// scheduleBgJobs schedules all jobs running in the background, e.g. cleaning yesterday's items from the list
func scheduleBgJobs() {
	// cron format: Minutes, Hours, day of Month, Month, Day of Week
	logger.Info("Scheduling jobs...")
	cfg.RunTime.BgJobs = cron.New()
	cfg.RunTime.BgJobs.Start()
	logger.Info("Jobs scheduled")
}

// startServer starts the preconfigured web server
func startServer() {
	logger.Infof("Listening on %v", cfg.RunTime.ListenAddr)
	cfg.RunTime.StartDate = date.GetNowUtc()
	if cfg.Server.UseTls {
		if err := server.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile); err != nil && err != http.ErrServerClosed {
			logger.Error("Error while starting https server", err)
			panic(err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error while starting http server", err)
			panic(err)
		}
	}
}

// cleanUp tries to clean up when the program is stopped
func cleanUp() {
	logger.Info("Cleaning up...")
	cfg.DeviceScan.DeviceScanRun = false
	cfg.RunTime.BgJobs.Stop()
	shutdownTime := time.Duration(cfg.Server.GracefulShutdownTime) * time.Second
	ctx, cancel = context.WithTimeout(context.Background(), shutdownTime)
	defer func() {
		logger.Info("Cleaned up")
		cancel()
	}()
}
