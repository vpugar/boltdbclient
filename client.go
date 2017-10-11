package boltdbclient

import (
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// Client represents a boltdbclient to the underlying BoltDB data store.
type Client struct {
	// Returns the current time.
	Now func() time.Time

	Db *bolt.DB

	config Config
}

// Creates new client
func NewClient(config Config) *Client {
	return &Client{
		config: config,
		Now:    time.Now,
	}
}

// Open opens and initializes the BoltDB database.
func (c *Client) Open() (string, error) {

	if err := os.MkdirAll(c.config.Dir, 0666); err != nil {
		return "", err
	}

	dbPath := path.Join(c.config.Dir, c.config.Filename)
	// Open database file.
	if dbConn, err := bolt.Open(dbPath, 0666, &bolt.Options{Timeout: 1 * time.Second}); err != nil {
		return dbPath, err
	} else {
		c.Db = dbConn
		return dbPath, nil
	}
}

// Close closes then underlying BoltDB database.
func (c *Client) Close() error {
	if c.Db != nil {
		return c.Db.Close()
	}
	return nil
}

type TransactionCallback func(*bolt.Tx) error

// Creation of initial entities (for example buckets)
func (client *Client) InitEntity(initEntityCallback TransactionCallback) error {
	if tx, err := client.Db.Begin(true); err != nil {
		return err
	} else {

		success := false
		defer DeferRollback(tx, &success)

		if err = initEntityCallback(tx); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		} else {
			success = true
			return nil
		}
	}
}

// Write transaction with callback
func (client *Client) WriteTransaction(write TransactionCallback) error {

	// Start write transaction.
	if tx, err := client.Db.Begin(true); err != nil {
		return errors.WithStack(err)
	} else {
		success := false
		defer DeferRollback(tx, &success)

		if err = write(tx); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return errors.WithStack(err)
		} else {
			success = true
			return nil
		}
	}
}

// Read transaction with callback
func (client *Client) ReadTransaction(read TransactionCallback) error {

	// Start read-only transaction.
	if tx, err := client.Db.Begin(false); err != nil {
		return errors.WithStack(err)
	} else {
		defer tx.Rollback()
		return read(tx)
	}
}

// Delete entry from bucket
func (client *Client) DeleteWithTransaction(bucket []byte, id string) error {

	// Start write transaction.
	if tx, err := client.Db.Begin(true); err != nil {
		return errors.WithStack(err)
	} else {
		success := false
		defer DeferRollback(tx, &success)

		b := tx.Bucket(bucket)

		if err := b.Delete([]byte(id)); err != nil {
			return errors.WithStack(err)
		}

		if err = tx.Commit(); err != nil {
			return errors.WithStack(err)
		} else {
			success = true
			return nil
		}
	}
}
