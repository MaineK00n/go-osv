package cmd

import (
	"github.com/MaineK00n/go-osv/db"
	"github.com/MaineK00n/go-osv/fetcher"
	"github.com/MaineK00n/go-osv/models"
	"github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

// dwfCmd represents the dwf command
var dwfCmd = &cobra.Command{
	Use:   "dwf",
	Short: "Fetch the CVE information from osv-vulnerabilities/DWF",
	Long:  `Fetch the CVE information from osv-vulnerabilities/DWF`,
	RunE:  fetchDWF,
}

func init() {
	fetchCmd.AddCommand(dwfCmd)
}

func fetchDWF(cmd *cobra.Command, args []string) (err error) {
	log15.Info("Initialize Database")
	driver, locked, err := db.NewDB(viper.GetString("dbtype"), viper.GetString("dbpath"), viper.GetBool("debug-sql"))
	if err != nil {
		if locked {
			log15.Error("Failed to initialize DB. Close DB connection before fetching", "err", err)
		}
		return err
	}

	fetchMeta, err := driver.GetFetchMeta()
	if err != nil {
		log15.Error("Failed to get FetchMeta from DB.", "err", err)
		return err
	}
	if fetchMeta.OutDated() {
		log15.Error("Failed to Insert CVEs into DB. SchemaVersion is old", "SchemaVersion", map[string]uint{"latest": models.LatestSchemaVersion, "DB": fetchMeta.SchemaVersion})
		return xerrors.New("Failed to Insert CVEs into DB. SchemaVersion is old")
	}

	log15.Info("Fetched all OSV Data from osv-vulnerabilities/DWF")
	osvJSONs, err := fetcher.FetchOSVs(models.DWF)
	if err != nil {
		log15.Error("Failed to Fetch OSV Data from osv-vulnerabilities/DWF.", "err", err)
		return err
	}

	log15.Info("Fetched", "OSVs", len(osvJSONs))

	log15.Info("Insert OSVs into DB", "db", driver.Name())
	if err := driver.InsertOSVs(models.DWF, osvJSONs); err != nil {
		log15.Error("Failed to insert.", "dbpath",
			viper.GetString("dbpath"), "err", err)
		return err
	}

	if err := driver.UpsertFetchMeta(fetchMeta); err != nil {
		log15.Error("Failed to upsert FetchMeta to DB.", "err", err)
		return err
	}

	return nil
}
