package cmd

import (
	"github.com/MaineK00n/go-osv/db"
	"github.com/MaineK00n/go-osv/models"
	"github.com/MaineK00n/go-osv/server"
	"github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start OSV HTTP server",
	Long:  `Start OSV HTTP server`,
	RunE:  executeServer,
}

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().String("bind", "127.0.0.1", "HTTP server bind to IP address (default: loop back interface")
	viper.BindPFlag("bind", serverCmd.PersistentFlags().Lookup("bind"))

	serverCmd.PersistentFlags().String("port", "1328", "HTTP server port number (default: 1328")
	viper.BindPFlag("port", serverCmd.PersistentFlags().Lookup("port"))
}

func executeServer(cmd *cobra.Command, args []string) (err error) {
	logDir := viper.GetString("log-dir")
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
		log15.Error("Failed to start server. SchemaVersion is old", "SchemaVersion", map[string]uint{"latest": models.LatestSchemaVersion, "DB": fetchMeta.SchemaVersion})
		return xerrors.New("Failed to start server. SchemaVersion is old")
	}

	log15.Info("Starting HTTP Server...")
	if err = server.Start(logDir, driver); err != nil {
		log15.Error("Failed to start server.", "err", err)
		return err
	}

	return nil
}
