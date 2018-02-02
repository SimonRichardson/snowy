package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/SimonRichardson/flagset"
	"github.com/SimonRichardson/gexec"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/trussle/fsys"
	"github.com/trussle/snowy/pkg/contents"
	"github.com/trussle/snowy/pkg/journals"
	"github.com/trussle/snowy/pkg/ledgers"
	"github.com/trussle/snowy/pkg/repository"
	"github.com/trussle/snowy/pkg/status"
	"github.com/trussle/snowy/pkg/store"
	"github.com/trussle/snowy/pkg/ui"
)

const (
	defaultFilesystem  = "remote"
	defaultPersistence = "real"

	defaultAWSEncryption           = false
	defaultAWSKMSKey               = ""
	defaultAWSServerSideEncryption = "aws:kmskey"
	defaultAWSID                   = ""
	defaultAWSSecret               = ""
	defaultAWSToken                = ""
	defaultAWSRegion               = "eu-west-1"
	defaultAWSBucket               = ""

	defaultDBHostname = "localhost"
	defaultDBPort     = 54321
	defaultDBUsername = "postgres"
	defaultDBPassword = "postgres"
	defaultDBName     = "postgres"
	defaultDBSSLMode  = "disable"

	defaultMetricsRegistration = true
	defaultUILocal             = false
)

func runDocuments(args []string) error {
	// flags for the documents command
	var (
		flags = flagset.NewFlagSet("documents", flag.ExitOnError)

		debug                   = flags.Bool("debug", false, "debug logging")
		apiAddr                 = flags.String("api", defaultAPIAddr, "listen address for query API")
		filesystem              = flags.String("filesystem", defaultFilesystem, "type of filesystem backing (local, remote, virtual, nop)")
		datastore               = flags.String("persistence", defaultPersistence, "type of persistence backing (real, virtual, nop)")
		awsEncryption           = flags.Bool("aws.encryption", defaultAWSEncryption, "AWS configuration encryption")
		awsKMSKey               = flags.String("aws.kmskey", defaultAWSKMSKey, "AWS configuration KMS Key")
		awsServerSideEncryption = flags.String("aws.sse", defaultAWSServerSideEncryption, "AWS configuration ServerSideEncryption")
		awsID                   = flags.String("aws.id", defaultAWSID, "AWS configuration id")
		awsSecret               = flags.String("aws.secret", defaultAWSSecret, "AWS configuration secret")
		awsToken                = flags.String("aws.token", defaultAWSToken, "AWS configuration token")
		awsRegion               = flags.String("aws.region", defaultAWSRegion, "AWS configuration region")
		awsBucket               = flags.String("aws.bucket", defaultAWSBucket, "AWS configuration bucket")
		dbHost                  = flags.String("db.hostname", defaultDBHostname, "Host name for connecting to the the datastore")
		dbPort                  = flags.Int("db.port", defaultDBPort, "Port for connecting to the the datastore")
		dbUsername              = flags.String("db.username", defaultDBUsername, "Username for connecting to the datastore")
		dbPassword              = flags.String("db.password", defaultDBPassword, "Password for connecting to the datastore")
		dbName                  = flags.String("db.name", defaultDBName, "Name of the database with in the datastore")
		dbSSLMode               = flags.String("db.sslmode", defaultDBSSLMode, "SSL mode for connecting to the datastore")
		metricsRegistration     = flags.Bool("metrics.registration", defaultMetricsRegistration, "Registration of metrics on launch")
		uiLocal                 = flags.Bool("ui.local", defaultUILocal, "Ignores embedded files and goes straight to the filesystem")
	)

	flags.Usage = usageFor(flags, "documents [flags]")
	if err := flags.Parse(args); err != nil {
		return nil
	}

	// Setup the logger.
	var logger log.Logger
	{
		logLevel := level.AllowInfo()
		if *debug {
			logLevel = level.AllowAll()
		}
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = level.NewFilter(logger, logLevel)
	}

	// Instrumentation
	connectedClients := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "snowy_documents",
		Name:      "connected_clients",
		Help:      "Number of currently connected clients by modality.",
	}, []string{"modality"})
	writerBytes := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "snowy_documents",
		Name:      "writer_bytes_written_total",
		Help:      "The total number of bytes written.",
	})
	writerRecords := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "snowy_documents",
		Name:      "writer_records_written_total",
		Help:      "The total number of records written.",
	})
	apiDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "snowy_documents",
		Name:      "api_request_duration_seconds",
		Help:      "API request duration in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path", "status_code"})

	if *metricsRegistration {
		prometheus.MustRegister(
			connectedClients,
			writerBytes,
			writerRecords,
			apiDuration,
		)
	}

	apiNetwork, apiAddress, err := parseAddr(*apiAddr, defaultAPIPort)
	if err != nil {
		return err
	}
	apiListener, err := net.Listen(apiNetwork, apiAddress)
	if err != nil {
		return err
	}
	level.Debug(logger).Log("API", fmt.Sprintf("%s://%s", apiNetwork, apiAddress))

	// Filesystem setup.
	remoteConfig, err := fsys.BuildConfig(
		fsys.WithEncryption(*awsEncryption),
		fsys.WithID(*awsID),
		fsys.WithSecret(*awsSecret),
		fsys.WithToken(*awsToken),
		fsys.WithKMSKey(*awsKMSKey),
		fsys.WithServerSideEncryption(*awsServerSideEncryption),
		fsys.WithRegion(*awsRegion),
		fsys.WithBucket(*awsBucket),
	)
	if err != nil {
		return errors.Wrap(err, "filesystem remote config")
	}

	fysConfig, err := fsys.Build(
		fsys.With(*filesystem),
		fsys.WithConfig(remoteConfig),
	)
	if err != nil {
		return errors.Wrap(err, "filesystem config")
	}

	fsys, err := fsys.New(fysConfig)
	if err != nil {
		return errors.Wrap(err, "filesystem")
	}

	// Persistence setup.
	realConfig, err := store.BuildConfig(
		store.WithHostPort(*dbHost, *dbPort),
		store.WithUsername(*dbUsername),
		store.WithPassword(*dbPassword),
		store.WithDBName(*dbName),
		store.WithSSLMode(*dbSSLMode),
	)
	if err != nil {
		return errors.Wrap(err, "store real config")
	}

	storeConfig, err := store.Build(
		store.With(*datastore),
		store.WithConfig(realConfig),
	)
	if err != nil {
		return errors.Wrap(err, "store config")
	}

	dataStore, err := store.New(storeConfig, log.With(logger, "component", "store"))
	if err != nil {
		return errors.Wrap(err, "store")
	}

	// Repository setup
	repository := repository.NewRealRepository(fsys, dataStore, log.With(logger, "component", "repository"))
	defer func() {
		if err := repository.Close(); err != nil {
			level.Error(logger).Log("err", err.Error())
		}
	}()

	// Execution group.
	g := gexec.NewGroup()
	gexec.Block(g)
	{
		// Store manages and maintains the underlying dataStore.
		g.Add(func() error {
			err := dataStore.Run()
			return err
		}, func(error) {
			dataStore.Stop()
		})
	}
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			dst := make(chan struct{})
			go func() {
				for {
					select {
					case <-dst:
						stats, err := repository.LedgerStatistics()
						if err != nil {
							level.Error(logger).Log("err", err)
							return
						}

						t := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
						fmt.Fprintf(t, "Timestamp \tTotal \t\n")
						fmt.Fprintf(t, "%s \t%d \t\n", time.Now().Format(time.RFC3339), stats.TotalLedgers)
						t.Flush()
					}
				}
			}()
			return interrupt(cancel, dst)
		}, func(error) {
			close(cancel)
		})
	}
	{
		g.Add(func() error {
			contentsAPI := contents.NewAPI(repository,
				log.With(logger, "component", "contents_api"),
				connectedClients.WithLabelValues("contents"),
				writerBytes, writerRecords,
				apiDuration,
			)
			defer contentsAPI.Close()

			mux := http.NewServeMux()
			mux.Handle("/ledgers/", http.StripPrefix("/ledgers", ledgers.NewAPI(repository,
				log.With(logger, "component", "ledgers_api"),
				connectedClients.WithLabelValues("ledgers"),
				apiDuration,
			)))
			mux.Handle("/contents/", http.StripPrefix("/contents", contentsAPI))
			mux.Handle("/journals/", http.StripPrefix("/journals", journals.NewAPI(repository,
				log.With(logger, "component", "journals_api"),
				connectedClients.WithLabelValues("journals"),
				writerBytes, writerRecords,
				apiDuration,
			)))
			mux.Handle("/status/", http.StripPrefix("/status", status.NewAPI(
				log.With(logger, "component", "status_api"),
				connectedClients.WithLabelValues("status"),
				apiDuration,
			)))
			mux.Handle("/ui/", ui.NewAPI(*uiLocal, log.With(logger, "component", "ui")))

			registerMetrics(mux)
			registerProfile(mux)

			return http.Serve(apiListener, mux)
		}, func(error) {
			apiListener.Close()
		})
	}
	gexec.Interrupt(g)
	return g.Run()
}
