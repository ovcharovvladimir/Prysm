// Package database provides several constructs including a simple in-memory database.
// This should not be used for production, but would be a helpful interim
// solution for development.
package database

import (
	"fmt"
	"path/filepath"

	"github.com/ovcharovvladimir/essentiaHybrid/essdb"
	sharedDB "github.com/ovcharovvladimir/Prysm/shared/database"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "db")

// ShardDB defines a service for the sharding system's persistent storage.
type ShardDB struct {
	inmemory bool
	dataDir  string
	name     string
	cache    int
	handles  int
	db       essdb.Database
}

// ShardDBConfig specifies configuration options for the db service.
type ShardDBConfig struct {
	DataDir  string
	Name     string
	InMemory bool
}

// NewShardDB initializes a shardDB.
func NewShardDB(config *ShardDBConfig) (*ShardDB, error) {
	// Uses default cache and handles values.
	// TODO: allow these arguments to be set based on cli context.
	shardDB := &ShardDB{
		name:    config.Name,
		dataDir: config.DataDir,
	}
	if config.InMemory {
		shardDB.inmemory = true
		shardDB.db = sharedDB.NewKVStore()
	} else {
		shardDB.inmemory = false
		shardDB.cache = 16
		shardDB.handles = 16
	}
	return shardDB, nil
}

// Start the shard DB service.
func (s *ShardDB) Start() {
	log.Info("Starting shardDB service")
	if !s.inmemory {
		db, err := essdb.NewLDBDatabase(filepath.Join(s.dataDir, s.name), s.cache, s.handles)
		if err != nil {
			log.Error(fmt.Sprintf("Could not start shard DB: %v", err))
			return
		}
		s.db = db
	}
}

// Stop the shard DB service gracefully.
func (s *ShardDB) Stop() error {
	log.Info("Stopping shardDB service")
	s.db.Close()
	return nil
}

// DB returns the attached essdb instance.
func (s *ShardDB) DB() essdb.Database {
	return s.db
}
