package db

import (
	"github.com/gocql/gocql"
)

// Db represents a connection to a db
type Db struct {
	session DbSession
}

// NewDb Gets a pointer to a db
func NewDb(hosts ...string) (*Db, error) {
	cluster := gocql.NewCluster(hosts...)

	var (
		session *gocql.Session
		err     error
	)

	if session, err = cluster.CreateSession(); err != nil {
		return nil, err
	}

	return &Db{
		session: &GoCqlSession{ref: session},
	}, nil
}

// Keyspace Retrieves a keyspace
func (db *Db) Keyspace(keyspace string) (*gocql.KeyspaceMetadata, error) {
	// We expose gocql types for now, we should wrap them in the future instead
	return db.session.KeyspaceMetadata(keyspace)
}

// Keyspaces Retrieves all the keyspace names
func (db *Db) Keyspaces() ([]string, error) {
	iter := db.session.ExecuteIterSimple("SELECT keyspace_name FROM system_schema.keyspaces", gocql.One)

	var keyspaces []string

	var name string
	for iter.Scan(&name) {
		keyspaces = append(keyspaces, name)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return keyspaces, nil
}
