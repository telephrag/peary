package strg

import (
	"context"
	"encoding/json"
	"fmt"
	"kubinka/bot_errors"
	"kubinka/config"
	"kubinka/models"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
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

func (b *BoltConn) Close() error {
	return b.db.Close()
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

func (b *BoltConn) WatchExpirations(ctx context.Context, ds *discordgo.Session) error {
	timeout := time.After(time.Second * config.DB_CHANGESTREAM_SLEEP_SECONDS)

	for {
		select {
		case <-timeout:
			tx, err := b.db.Begin(true)
			if err != nil {
				return fmt.Errorf("bolt: failed to initiate transaction: %w", err)
			}

			bkt := tx.Bucket([]byte(b.domain))
			c := bkt.Cursor()
			now := time.Now()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				p := models.Player{}
				if err := json.Unmarshal(v, &p); err != nil {
					return fmt.Errorf("bolt: failed to unmarshal record into models.Player: %w", err)
				}

				if p.Expire.Before(now) {
					/* ErrIncompatibleValue may occur if you have nested bucket
					which may become the case in the future. */
					err := ds.GuildMemberRoleRemove(config.BOT_GUILD_ID, p.DiscordID, config.BOT_ROLE_ID)
					if err != nil {
						return fmt.Errorf(
							"discord: failed to remove role from player %s: %w",
							p.DiscordID,
							err,
						)
					}
					err = c.Delete()
					if err != nil {
						return fmt.Errorf("bolt: failed to delete player record at cursor: %w", err)
					}

					log.Printf("%s deployment time expired, removed from db\n", p.DiscordID)
				}
			}
			if err := tx.Commit(); err != nil {
				tx.Rollback()
				return fmt.Errorf("bolt: failed to commit transaction: %w", err)
			}
			timeout = time.After(time.Second * config.DB_CHANGESTREAM_SLEEP_SECONDS)

		case <-ctx.Done():
			return fmt.Errorf(bot_errors.ErrSomewhereElse)
		}
	}
}
