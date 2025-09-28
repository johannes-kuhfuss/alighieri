package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/alighieri/config"
	"github.com/johannes-kuhfuss/alighieri/repositories"
	"github.com/stretchr/testify/assert"
)

var (
	repo     repositories.DefaultDeviceRepository
	uh       StatsUiHandler
	cfg      config.AppConfig
	router   *gin.Engine
	recorder *httptest.ResponseRecorder
)

func setupUiTest() func() {
	config.InitConfig("", &cfg)
	repo = repositories.NewDeviceRepository(&cfg)
	uh = NewStatsUiHandler(&cfg, &repo)
	router = gin.Default()
	router.LoadHTMLGlob("../templates/*.tmpl")
	recorder = httptest.NewRecorder()
	return func() {
		router = nil
	}
}

func TestStatusPageReturnsStatus(t *testing.T) {
	teardown := setupUiTest()
	defer teardown()
	router.GET("/", uh.StatusPage)
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	router.ServeHTTP(recorder, request)
	res := recorder.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	containsTitle := strings.Contains(string(data), "<title>Status</title>")

	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)
	assert.True(t, containsTitle)
}

func TestAboutPageReturnsAbout(t *testing.T) {
	teardown := setupUiTest()
	defer teardown()
	router.GET("/about", uh.AboutPage)
	request := httptest.NewRequest(http.MethodGet, "/about", nil)

	router.ServeHTTP(recorder, request)
	res := recorder.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	containsTitle := strings.Contains(string(data), "<title>About</title>")

	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)
	assert.True(t, containsTitle)
}

func TestDeviceListPageReturnsDeviceListPage(t *testing.T) {
	teardown := setupUiTest()
	defer teardown()
	router.GET("/devicelist", uh.DeviceListPage)
	request := httptest.NewRequest(http.MethodGet, "/devicelist", nil)

	router.ServeHTTP(recorder, request)
	res := recorder.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	containsTitle := strings.Contains(string(data), "<title>Device List</title>")

	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)
	assert.True(t, containsTitle)
}

func TestLogsPageReturnsLogs(t *testing.T) {
	teardown := setupUiTest()
	defer teardown()
	router.GET("/logs", uh.LogsPage)
	request := httptest.NewRequest(http.MethodGet, "/logs", nil)

	router.ServeHTTP(recorder, request)
	res := recorder.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	containsTitle := strings.Contains(string(data), "<title>Logs</title>")

	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)
	assert.True(t, containsTitle)
}
