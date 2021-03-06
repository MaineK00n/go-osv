package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/MaineK00n/go-osv/db"
	"github.com/MaineK00n/go-osv/models"
	"github.com/inconshreveable/log15"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// Start starts OSV HTTP Server.
func Start(logDir string, driver db.DB) error {
	e := echo.New()
	e.Debug = viper.GetBool("debug")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// setup access logger
	logPath := filepath.Join(logDir, "access.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if _, err := os.Create(logPath); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: f,
	}))

	// Routes
	e.GET("/health", health())
	e.GET("/ids/:id", getOSVbyID(driver))
	e.GET("/:type/ids/:id", getOSVbyID(driver))
	e.GET("/pkgs/:name", getOSVbyPackageName(driver))
	e.GET("/:type/pkgs/:name", getOSVbyPackageName(driver))

	bindURL := fmt.Sprintf("%s:%s", viper.GetString("bind"), viper.GetString("port"))
	log15.Info("Listening", "URL", bindURL)

	return e.Start(bindURL)
}

// Handler
func health() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}
}

// Handler
func getOSVbyID(driver db.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		osvType := c.Param("type")
		log15.Debug("Params", "id", id, "osvType", osvType)
		osvDetail, err := driver.GetOSVbyID(id, osvType)
		if err != nil {
			log15.Error("Failed to GetOSVbyID.", "err", err)
			return c.JSON(http.StatusInternalServerError, []models.OSV{})
		}

		return c.JSON(http.StatusOK, osvDetail)
	}
}

// Handler
func getOSVbyPackageName(driver db.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		osvType := c.Param("type")
		log15.Debug("Params", "name", name, "osvType", osvType)
		osvDetail, err := driver.GetOSVbyPackageName(name, osvType)
		if err != nil {
			log15.Error("Failed to GetOSVbyPackageName.", "err", err)
			return c.JSON(http.StatusInternalServerError, []models.OSV{})
		}
		return c.JSON(http.StatusOK, osvDetail)
	}
}
