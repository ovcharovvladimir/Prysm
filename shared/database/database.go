// Package database defines a beacon chain DB service that can be
// initialized with either a persistent db, or an in-memory kv-store.
package database

import (
	"fmt"
	"path/filepath"

	"github.com/mattn/go-colorable"
	"github.com/ovcharovvladimir/essentiaHybrid/essdb"
	"github.com/ovcharovvladimir/essentiaHybrid/log"
	"gopkg.in/urfave/cli.v1"
)

// DB defines a service for the beacon chain system's persistent storage.
type DB struct {
	_db essdb.Database
}

// DBConfig specifies configuration options for the db service.
type DBConfig struct {
	DataDir  string
	Name     string
	InMemory bool
}

// NewDB initializes a beaconDB instance.
func NewDB(config *DBConfig) (*DB, error) {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(3), log.StreamHandler(colorable.NewColorableStdout(), log.TerminalFormat(true))))
	// Uses default cache and handles values.
	db := &DB{}
	if config.InMemory {
		db._db = NewKVStore()
		return db, nil
	}

	_db, err := essdb.NewLDBDatabase(filepath.Join(config.DataDir, config.Name), 16, 16)
	if err != nil {
		log.Error("Could not start beaconDB", "err", err)
		return nil, err
	}
	db._db = _db

	return db, nil
}

// Close closes the database
func (b *DB) Close() {
	b._db.Close()
}

// DB returns the attached ethdb instance.
func (b *DB) DB() essdb.Database {
	return b._db
}
