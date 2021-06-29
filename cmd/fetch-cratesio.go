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

// cratesioCmd represents the cratesio command
var cratesioCmd = &cobra.Command{
	Use:   "crates.io",
	Short: "Fetch the CVE information from osv-vulnerabilities/crates.io",
	Long:  `Fetch the CVE information from osv-vulnerabilities/crates.io`,
	RunE:  fetchCratesio,
}

func init() {
	fetchCmd.AddCommand(cratesioCmd)
}

func fetchCratesio(cmd *cobra.Command, args []string) (err error) {
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

	log15.Info("Fetched all OSV Data from osv-vulnerabilities/crates.io")
	osvs, err := fetcher.FetchOSVDetails(models.CratesIOType)
	if err != nil {
		return err
	}

	log15.Info("Fetched", "OSVs", len(osvs))

	if err := driver.UpsertFetchMeta(fetchMeta); err != nil {
		log15.Error("Failed to upsert FetchMeta to DB.", "err", err)
		return err
	}

	return nil
}
