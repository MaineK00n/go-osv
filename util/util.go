package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/inconshreveable/log15"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
)

// GetDefaultLogDir returns default log directory
func GetDefaultLogDir() string {
	defaultLogDir := "/var/log/go-osv"
	if runtime.GOOS == "windows" {
		defaultLogDir = filepath.Join(os.Getenv("APPDATA"), "go-osv")
	}
	return defaultLogDir
}

// SetLogger set logger
func SetLogger(logDir string, debug, logJSON bool) {
	stderrHundler := log15.StderrHandler
	logFormat := log15.LogfmtFormat()
	if logJSON {
		logFormat = log15.JsonFormatEx(false, true)
		stderrHundler = log15.StreamHandler(os.Stderr, logFormat)
	}

	lvlHundler := log15.LvlFilterHandler(log15.LvlInfo, stderrHundler)
	if debug {
		lvlHundler = log15.LvlFilterHandler(log15.LvlDebug, stderrHundler)
	}

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.Mkdir(logDir, 0700); err != nil {
			log15.Error("Failed to create log directory", "err", err)
		}
	}
	var hundler log15.Handler
	if _, err := os.Stat(logDir); err == nil {
		logPath := filepath.Join(logDir, "go-osv.log")
		if _, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
			log15.Error("Failed to create a log file", "err", err)
			hundler = lvlHundler
		} else {
			hundler = log15.MultiHandler(
				log15.Must.FileHandler(logPath, logFormat),
				lvlHundler,
			)
		}
	} else {
		hundler = lvlHundler
	}
	log15.Root().SetHandler(hundler)
}

// FetchURL returns HTTP response body
func FetchURL(url string) ([]byte, error) {
	var errs []error
	httpProxy := viper.GetString("http-proxy")

	req := gorequest.New().Proxy(httpProxy).Get(url)
	resp, respBody, errs := req.End()
	if len(errs) > 0 || resp == nil {
		return nil, fmt.Errorf("HTTP error. errs: %v, url: %s", errs, url)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error. errs: %v, status code: %d, url: %s", errs, resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(strings.NewReader(respBody))
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll error. err: %w", err)
	}

	return body, nil
}
