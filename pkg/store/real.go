package store

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/lib/pq"

	"github.com/trussle/snowy/pkg/uuid"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	// pq is required here, even though it's not used, so that it gets injected
	// into database/sql - magic!
	_ "github.com/lib/pq"
)

const (
	getSelectQuery = "SELECT id, name, resource_id, author_id, tags, created_on, deleted_on FROM documents WHERE resource_id = $1;"
	getInsertQuery = "INSERT INTO documents (name, resource_id, author_id, tags, created_on, deleted_on) VALUES ($1, $2, $3, $4, $5, $6);"
)

// RealConfig holds the options for connecting to the DB
type RealConfig struct {
	Host               string
	Port               int
	Username, Password string
	DBName             string
	SSLMode            string
}

type realStore struct {
	config *RealConfig
	db     *sql.DB
	stop   chan chan struct{}
	logger log.Logger
}

// NewRealStore yields a real data store.
func NewRealStore(config *RealConfig, logger log.Logger) Store {
	return &realStore{
		config: config,
		logger: logger,
	}
}

func (r *realStore) Get(resource uuid.UUID) (Entity, error) {
	var (
		entity Entity
		row    = r.db.QueryRow(getSelectQuery, resource.String())

		resourceID, authorID string
	)
	err := row.Scan(
		&entity.ID,
		&entity.Name,
		&resourceID,
		&authorID,
		pq.Array(&entity.Tags),
		&entity.CreatedOn,
		&entity.DeletedOn,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity, errNotFound{err}
		}
		return entity, err
	}

	// We have to manually extract the UUID, as database/sql doesn't provide this
	// for us.
	if entity.ResourceID, err = uuid.Parse(resourceID); err != nil {
		return entity, err
	}
	if entity.AuthorID, err = uuid.Parse(authorID); err != nil {
		return entity, err
	}

	return entity, nil
}

func (r *realStore) Put(entity Entity) error {
	return r.Transaction(func(txn *sql.Tx) error {
		stmt, err := txn.Prepare(getInsertQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		if _, err = stmt.Exec(
			entity.Name,
			entity.ResourceID.String(),
			entity.AuthorID.String(),
			pq.Array(entity.Tags),
			entity.CreatedOn,
			entity.DeletedOn,
		); err != nil {
			return errors.Wrap(err, "unable to exec statement")
		}

		return nil
	})
}

func (r *realStore) Transaction(fn func(*sql.Tx) error) (err error) {
	var txn *sql.Tx
	txn, err = r.db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			txn.Rollback()
		} else {
			if err = txn.Commit(); err != nil {
				err = errors.Wrap(err, "unable to commit statement")
			}
		}
	}()

	err = fn(txn)
	return
}

// Run the store
func (r *realStore) Run() error {
	var err error
	r.db, err = sql.Open("postgres", ConnectionString(r.config))
	if err != nil {
		level.Error(r.logger).Log("err", err)
		return err
	}

	// Sync db check to validate conn-opts
	if err = r.db.Ping(); err != nil {
		level.Error(r.logger).Log("err", err)
		return err
	}

	for {
		select {
		case c := <-r.stop:
			err := r.db.Close()
			close(c)
			return err
		}
	}
}

// Stop the store
func (r *realStore) Stop() {
	c := make(chan struct{})
	r.stop <- c
	<-c
}

// RealOption defines a option for generating a RealConfig
type RealOption func(*RealConfig) error

// BuildConfig ingests configuration options to then yield a RealConfig, and return an
// error if it fails during configuring.
func BuildConfig(opts ...RealOption) (*RealConfig, error) {
	var config RealConfig
	for _, opt := range opts {
		err := opt(&config)
		if err != nil {
			return nil, err
		}
	}
	return &config, nil
}

// WithHostPort adds a host and port option to the configuration
func WithHostPort(host string, port int) RealOption {
	return func(config *RealConfig) error {
		config.Host = host
		config.Port = port
		return nil
	}
}

// WithUsername adds a username option to the configuration
func WithUsername(username string) RealOption {
	return func(config *RealConfig) error {
		config.Username = username
		return nil
	}
}

// WithPassword adds a password option to the configuration
func WithPassword(password string) RealOption {
	return func(config *RealConfig) error {
		config.Password = password
		return nil
	}
}

// WithDBName adds a db name option to the configuration
func WithDBName(dbName string) RealOption {
	return func(config *RealConfig) error {
		config.DBName = dbName
		return nil
	}
}

// WithSSLMode adds a db name option to the configuration
func WithSSLMode(sslMode string) RealOption {
	return func(config *RealConfig) error {
		switch sslMode {
		case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
			break
		default:
			return errors.Errorf("unexpected ssl mode: %s", sslMode)
		}
		config.SSLMode = sslMode
		return nil
	}
}

// ConnectionString consumes a configuration file and returns a connection
// string to the database.
func ConnectionString(config *RealConfig) string {
	var opts []string

	for _, v := range []struct {
		key, value string
	}{
		{"host", config.Host},
		{"port", strconv.Itoa(config.Port)},
		{"user", config.Username},
		{"password", config.Password},
		{"dbname", config.DBName},
		{"sslmode", config.SSLMode},
	} {
		if s := strings.TrimSpace(v.value); s != "" {
			opts = append(opts, fmt.Sprintf("%s=%s", v.key, v.value))
		}
	}

	return strings.Join(opts, " ")
}
