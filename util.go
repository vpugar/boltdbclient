package boltdbclient

import (
	"encoding/binary"

	"github.com/boltdb/bolt"
)

func DeferRollback(tx *bolt.Tx, success *bool) {
	if !(*success) {
		tx.Rollback()
	}
}

// Find nested bucked according to path.
// In case there is no some bucket in path, then return last not found name of bucket.
func FindBucket(b *bolt.Bucket, path []string) (*bolt.Bucket, string) {
	lb := b
	for _, part := range path {
		if lb = lb.Bucket([]byte(part)); lb == nil {
			return lb, part
		}
	}
	return lb, ""
}

func I2B(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
