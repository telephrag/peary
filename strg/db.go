package strg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"peary/config"
	"peary/errlist"
	"peary/models"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.etcd.io/bbolt"
)

type BoltConn struct {
	db     *bbolt.DB
	domain string
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

func (b *BoltConn) GetPlayerIDs() []string {
	var ids []string
	b.db.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(b.domain))
		ids = make([]string, bkt.Stats().KeyN)

		c := bkt.Cursor()
		i := 0
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			ids[i] = string(k)
			i++
		}
		return nil
	})

	return ids
}

func (b *BoltConn) Insert(p *models.Player) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(b.domain))

		buf, err := json.Marshal(p)
		if err != nil {
			return errlist.New(fmt.Errorf("failed to marshal player: %w", err)).
				Set("session", p.DiscordID).
				Set("event", errlist.DBInsert)
		}

		err = bkt.Put([]byte(p.DiscordID), buf)
		if err != nil {
			return errlist.New(fmt.Errorf("bolt: failed to put player into bucket %s: %w", b.domain, err)).
				Set("session", p.DiscordID).
				Set("event", errlist.DBInsert)
		}

		log.Print(errlist.New(nil).
			Set("session", p.DiscordID).
			Set("event", errlist.DBInsert),
		)
		return nil
	})
}

func (b *BoltConn) Delete(key string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(b.domain))

		if err := bkt.Delete([]byte(key)); err != nil {
			return errlist.New(fmt.Errorf("failed to delete value at key: %w", err)).
				Set("session", key).
				Set("event", errlist.DBDelete)
		}
		return nil
	})
}

func (b *BoltConn) RemoveExpired(ds *discordgo.Session) error {
	tx, err := b.db.Begin(true)
	if err != nil {
		return errlist.New(fmt.Errorf("failed to initiate transaction")).
			Set("event", errlist.DBChangeStream)
	}

	bkt := tx.Bucket([]byte(b.domain))
	c := bkt.Cursor()
	now := time.Now()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		p := models.Player{}
		if err := json.Unmarshal(v, &p); err != nil {
			return errlist.New(err).
				Set("event", errlist.DBChangeStream)
		}

		if p.Expire.Before(now) {
			/* ErrIncompatibleValue may occur if you have nested bucket
			which may become the case in the future. */
			err := ds.GuildMemberRoleRemove(
				ds.State.Ready.Guilds[0].ID,
				p.DiscordID,
				config.BOT_ROLE_ID)
			if err != nil {
				return errlist.New(errlist.ErrFailedTakeRole).
					Set("session", p.DiscordID).
					Set("event", errlist.DBChangeStreamExpire)
			}
			err = c.Delete()
			if err != nil {
				return errlist.New(fmt.Errorf("failed to delete player record at cursor: %w", err)).
					Set("event", errlist.DBChangeStreamExpire)
			}

			log.Print(errlist.New(nil).
				Set("session", p.DiscordID).
				Set("event", errlist.DBChangeStreamExpire),
			)
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errlist.New(fmt.Errorf("failed to commit transaction: %w", err)).
			Set("event", errlist.DBChangeStream)
	}

	return nil
}

func (b *BoltConn) WatchExpirations(ctx context.Context, ds *discordgo.Session) error {
	timeout := time.After(time.Second * 0) // expired records are deleted immediately on startup

	for {
		select {
		case <-timeout:
			if err := b.RemoveExpired(ds); err != nil {
				return err
			}
			timeout = time.After(time.Second * config.DB_CHANGESTREAM_SLEEP_SECONDS)
		case <-ctx.Done():
			return errlist.New(ctx.Err())
		}
	}
}
