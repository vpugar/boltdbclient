package boltdbclient_test

import (
	"testing"

	"os"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/vpugar/boltdbclient"
)

const DB_NAME = "bolt_test.db"
const BUCKET = "BUCK"
const NOT_BUCKET = "NOT_BUCK"

var BUCKET_BYTES = []byte(BUCKET)
var NOT_BUCKET_BYTES = []byte(NOT_BUCKET)

func run(m *testing.M) int {
	defer os.Remove(DB_NAME)
	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func TestOpenClose(t *testing.T) {
	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: DB_NAME,
	})
	if _, err := boltClient.Open(); err != nil {
		t.Fatal(err)
	}
	if err := boltClient.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestCloseWithoutOpen(t *testing.T) {
	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: DB_NAME,
	})
	if err := boltClient.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestInitEntity(t *testing.T) {
	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: DB_NAME,
	})
	if _, err := boltClient.Open(); err != nil {
		t.Fatal(err)
	}
	defer boltClient.Close()

	if err := boltClient.InitEntity(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(BUCKET_BYTES); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if err := boltClient.ReadTransaction(func(tx *bolt.Tx) error {
		if tx.Bucket(BUCKET_BYTES) == nil {
			return errors.New("No bucket")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestWriteTransaction(t *testing.T) {
	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: DB_NAME,
	})
	if _, err := boltClient.Open(); err != nil {
		t.Fatal(err)
	}
	defer boltClient.Close()

	if err := boltClient.WriteTransaction(func(tx *bolt.Tx) error {
		if b, err := tx.CreateBucketIfNotExists(BUCKET_BYTES); err != nil {
			return err
		} else {
			b.Put([]byte("test"), []byte("val"))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if err := boltClient.ReadTransaction(func(tx *bolt.Tx) error {
		if b := tx.Bucket(BUCKET_BYTES); b == nil {
			return errors.New("No bucket")
		} else {
			v := b.Get([]byte("test"))
			if string(v) != "val" {
				return errors.New("Val not stored")
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestWriteTransactionThatFails(t *testing.T) {
	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: DB_NAME,
	})
	if _, err := boltClient.Open(); err != nil {
		t.Fatal(err)
	}
	defer boltClient.Close()

	if err := boltClient.WriteTransaction(func(tx *bolt.Tx) error {
		if b, err := tx.CreateBucketIfNotExists(NOT_BUCKET_BYTES); err != nil {
			return err
		} else {
			b.Put([]byte("test"), []byte("val"))
		}
		return errors.New("Failed")
	}); err == nil {
		t.Fatal("Expected error")
	}

	if err := boltClient.ReadTransaction(func(tx *bolt.Tx) error {
		if b := tx.Bucket(NOT_BUCKET_BYTES); b == nil {
			// expected that there is no bucket
			return nil
		}
		return errors.New("Failed")
	}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEntryFromBucket(t *testing.T) {

	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: DB_NAME,
	})
	if _, err := boltClient.Open(); err != nil {
		t.Fatal(err)
	}
	defer boltClient.Close()

	if err := boltClient.WriteTransaction(func(tx *bolt.Tx) error {
		if b, err := tx.CreateBucketIfNotExists(BUCKET_BYTES); err != nil {
			return err
		} else {
			b.Put([]byte("test"), []byte("val"))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	boltClient.DeleteWithTransaction(BUCKET_BYTES, "test")

	if err := boltClient.ReadTransaction(func(tx *bolt.Tx) error {
		if b := tx.Bucket(BUCKET_BYTES); b != nil {
			v := b.Get([]byte("test"))
			if v != nil {
				return errors.New("Val still exists")
			}
			return nil
		} else {
			return errors.New("No bucket")
		}
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFindBucket(t *testing.T) {

	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: DB_NAME,
	})
	if _, err := boltClient.Open(); err != nil {
		t.Fatal(err)
	}
	defer boltClient.Close()

	if err := boltClient.WriteTransaction(func(tx *bolt.Tx) error {
		if b1, err := tx.CreateBucketIfNotExists([]byte("b1")); err != nil {
			return err
		} else {
			if b2, err := b1.CreateBucketIfNotExists([]byte("b2")); err != nil {
				return err
			} else {
				if _, err := b2.CreateBucketIfNotExists([]byte("b3")); err != nil {
					return err
				}
			}

			if b3, part := boltdbclient.FindBucket(b1, []string{"b2", "b3"}); b3 == nil {
				return errors.New("No bucket b1/b2/b3")
			} else if part != "" {
				return errors.New(part)
			}

			if b3, part := boltdbclient.FindBucket(b1, []string{"b2", "b4"}); b3 == nil {
				// No bucket part
				if part != "b4" {
					return errors.New("Wrong no bucket string " + part)
				}
			} else {
				return errors.New("Last found bucket must be empty " + part)
			}
		}

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
