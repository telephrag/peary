package strg

import (
	"encoding/json"
	"fmt"
	"kubinka/models"

	"go.etcd.io/bbolt"
)

type BoltConn struct {
	db     *bbolt.DB
	domain string
}

func (b *BoltConn) init(db *bbolt.DB, domain string) {
	b.db = db
	b.domain = domain
}

func Connect(dbName, domain string) (*BoltConn, error) {
	db, err := bbolt.Open(dbName+".db", 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %w", dbName, err)
	}

	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(domain))
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", domain, err)
		}
		return nil
	})

	return &BoltConn{
		db:     db,
		domain: domain,
	}, nil
}

func (b *BoltConn) Insert(p *models.Player) error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(b.domain))

		buf, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("failed to marshal player %v: %w", p, err)
		}

		err = bkt.Put([]byte(p.DiscordID), buf)
		if err != nil {
			return fmt.Errorf("failed to put player %v into bucket %s: %w", p, b.domain, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed to add %v to db: %w", p, err)
	}

	return nil
}

func (b *BoltConn) Delete(key string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(b.domain))

		if err := bkt.Delete([]byte(key)); err != nil {
			return fmt.Errorf("failed to delete value at key %s: %w", key, err)
		}

		return nil
	})
}
