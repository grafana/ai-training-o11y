package main

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	collector_version "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"

	app "github.com/grafana/ai-o11y/metadata-service/app"
	db "github.com/grafana/ai-o11y/metadata-service/internal"
)

// Version is set via build flag -ldflags -X main.Version
var (
	Version  string
	Branch   string
	Revision string
)

func init() {
	version.Version = Version
	version.Branch = Branch
	version.Revision = Revision
	prometheus.MustRegister(collector_version.NewCollector("metadata_service"))
}

// initialise mux router

func main() {
	os.Exit(run())
}

// run is the main entry point for the metadata-service.
// It initializes the database connection, creates the server and router.
// It also returns the exit code.
func run() int {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default("0.0.0.0").String()
		listenPort = kingpin.Flag(
			"web.listen-port",
			"Port on which to expose metrics and web interface.",
		).Default("4032").Int()
		databaseAddress = kingpin.Flag(
			"database-address",
			"Database connection string.",
		).Default("file:aitraining.db?mode=memory&cache=shared").String()
		databaseType = kingpin.Flag(
			"database-type",
			"Database type.",
		).Default(db.SQLite).Enum(db.SQLite, db.MySQL)
		constTenant = kingpin.Flag(
			"const-tenant",
			"A constant tenant to add to every request. Should only be used in development.",
		).String()
		// tenantOverridesFile = kingpin.Flag(
		// 	"tenant-overrides-file",
		// 	"Path to YAML file containing overrides per tenant.",
		// ).ExistingFile()
	)

	// Allow configuration to be specified via environment variables.
	kingpin.CommandLine.DefaultEnvars()

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("metadata-service"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	_, err := app.New(listenAddress, listenPort, databaseAddress, databaseType, constTenant, promlogConfig)
	if err != nil {
		return 1
	}

	return 0
}
