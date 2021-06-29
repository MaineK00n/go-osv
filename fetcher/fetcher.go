package fetcher

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/MaineK00n/go-osv/models"
	"github.com/MaineK00n/go-osv/util"
	"github.com/inconshreveable/log15"
	"golang.org/x/xerrors"
)

func FetchOSVDetails(fetchOSVType models.OSVType) (osvs []models.OSVJSON, err error) {
	url := fmt.Sprintf("https://osv-vulnerabilities.storage.googleapis.com/%s/all.zip", fetchOSVType)
	body, err := util.FetchURL(url)
	if err != nil {
		return []models.OSVJSON{}, xerrors.Errorf("Failed to fetch OSV data from url: %s. err: %w", url, err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return []models.OSVJSON{}, xerrors.Errorf("Failed to open zip.NewReader(). err: %w", err)
	}

	for _, zipFile := range zipReader.File {
		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			log15.Error("err", err)
			continue
		}

		var osv models.OSVJSON
		if err := json.Unmarshal(unzippedFileBytes, &osv); err != nil {
			return []models.OSVJSON{}, xerrors.Errorf("Failed to json.Unmarshal. err: %w", err)
		}

		osvs = append(osvs, osv)
	}

	return osvs, nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
