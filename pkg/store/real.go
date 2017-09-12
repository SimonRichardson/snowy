package store

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/trussle/snowy/pkg/uuid"
)

const (
	defaultSelectQuery = `SELECT id, 
	parent_id, 
	name, 
	resource_id, 
	resource_address, 
	resource_size, 
	resource_content_type, 
	author_id, 
	tags, 
	created_on, 
	deleted_on 
FROM   ledgers 
WHERE  resource_id = $1 
ORDER  BY created_on DESC, 
		 resource_address DESC;`
	defaultInsertQuery = `INSERT INTO ledgers 
	(parent_id, 
	 name, 
	 resource_id, 
	 resource_address, 
	 resource_size, 
	 resource_content_type, 
	 author_id, 
	 tags, 
	 created_on, 
	 deleted_on) 
VALUES      ($1, 
	 $2, 
	 $3, 
	 $4, 
	 $5, 
	 $6, 
	 $7, 
	 $8, 
	 $9, 
	 $10);`
	defaultSelectQueryTags = `SELECT id, 
	parent_id, 
	name, 
	resource_id, 
	resource_address, 
	resource_size, 
	resource_content_type, 
	author_id, 
	tags, 
	created_on, 
	deleted_on 
FROM   ledgers 
WHERE  resource_id = $1 
	AND tags && $2 
ORDER  BY created_on DESC, 
		 resource_address DESC;`
	defaultSelectQueryTagsAuthorID = `SELECT id, 
	parent_id, 
	name, 
	resource_id, 
	resource_address, 
	resource_size, 
	resource_content_type, 
	author_id, 
	tags, 
	created_on, 
	deleted_on 
FROM   ledgers 
WHERE  resource_id = $1 
	AND author_id = $2 
	AND tags && $3 
ORDER  BY created_on DESC, 
		 resource_address DESC;`
	defaultDropQuery = `TRUNCATE TABLE ledgers;`
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
	ticker *time.Ticker
}

// NewRealStore yields a real data store.
func NewRealStore(config *RealConfig, logger log.Logger) Store {
	return &realStore{
		config: config,
		stop:   make(chan chan struct{}),
		logger: logger,
		ticker: time.NewTicker(time.Minute),
	}
}

func (r *realStore) Get(resource uuid.UUID, query Query) (Entity, error) {
	var (
		entity          Entity
		statement, args = buildSQLFromQuery(resource, query)
		row             = r.db.QueryRow(statement, args...)

		id, parentID, resourceID string
	)
	err := row.Scan(
		&id,
		&parentID,
		&entity.Name,
		&resourceID,
		&entity.ResourceAddress,
		&entity.ResourceSize,
		&entity.ResourceContentType,
		&entity.AuthorID,
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
	if entity.ID, err = uuid.Parse(id); err != nil {
		return entity, err
	}
	if entity.ParentID, err = uuid.Parse(parentID); err != nil {
		return entity, err
	}
	if entity.ResourceID, err = uuid.Parse(resourceID); err != nil {
		return entity, err
	}

	return entity, nil
}

func (r *realStore) Insert(entity Entity) error {
	return r.Transaction(func(txn *sql.Tx) error {
		stmt, err := txn.Prepare(defaultInsertQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		// Normalize the tags of the entity
		tags := sortTags(entity.Tags)

		if _, err = stmt.Exec(
			entity.ParentID.String(),
			entity.Name,
			entity.ResourceID.String(),
			entity.ResourceAddress,
			entity.ResourceSize,
			entity.ResourceContentType,
			entity.AuthorID,
			pq.Array(tags),
			entity.CreatedOn,
			entity.DeletedOn,
		); err != nil {
			return errors.Wrap(err, "unable to exec statement")
		}

		return nil
	})
}

func (r *realStore) GetMultiple(resource uuid.UUID, query Query) ([]Entity, error) {
	var (
		statement, args = buildSQLFromQuery(resource, query)
		rows, err       = r.db.Query(statement, args...)
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []Entity
	for rows.Next() {
		var (
			entity Entity

			id, parentID, resourceID string
		)
		err := rows.Scan(
			&id,
			&parentID,
			&entity.Name,
			&resourceID,
			&entity.ResourceAddress,
			&entity.ResourceSize,
			&entity.ResourceContentType,
			&entity.AuthorID,
			pq.Array(&entity.Tags),
			&entity.CreatedOn,
			&entity.DeletedOn,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errNotFound{err}
			}
			return nil, err
		}

		// We have to manually extract the UUID, as database/sql doesn't provide this
		// for us.
		if entity.ID, err = uuid.Parse(id); err != nil {
			return nil, err
		}
		if entity.ParentID, err = uuid.Parse(parentID); err != nil {
			return nil, err
		}
		if entity.ResourceID, err = uuid.Parse(resourceID); err != nil {
			return nil, err
		}

		res = append(res, entity)
	}

	return res, rows.Err()
}

func (r *realStore) Transaction(fn func(*sql.Tx) error) (err error) {
	if r.db == nil {
		err = errors.New("db not found")
		return
	}

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

// Drop removes all of the stored ledgers
func (r *realStore) Drop() error {
	if r.db == nil {
		return errors.New("db not found")
	}

	_, err := r.db.Exec(defaultDropQuery)
	level.Error(r.logger).Log("err", err)
	return err
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

	// Make sure that we spin and ping the db, to make sure we're still current.
	for {
		select {
		case <-r.ticker.C:
			if err = r.db.Ping(); err != nil {
				level.Error(r.logger).Log("err", err)
				// TODO: reconnect to the db...
			}

		case c := <-r.stop:
			// First stop the ticker
			r.ticker.Stop()

			// Shut the db down.
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

func buildSQLFromQuery(resourceID uuid.UUID, query Query) (string, []interface{}) {
	numTags, authorID := len(query.Tags), query.AuthorID

	if numTags == 0 && (authorID == nil || *authorID == "") {
		return defaultSelectQuery, []interface{}{resourceID.String()}
	}

	if numTags > 0 && (authorID == nil || *authorID == "") {
		return defaultSelectQueryTags, []interface{}{
			resourceID.String(),
			pq.Array(query.Tags),
		}
	}

	return defaultSelectQueryTagsAuthorID, []interface{}{
		resourceID.String(),
		*authorID,
		pq.Array(query.Tags),
	}
}

func sortTags(tags []string) []string {
	res := make([]string, len(tags))
	copy(res, tags)
	sort.Strings(res)
	return res
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
