package db

import (
	"context"
	"database/sql"
	upperdb "github.com/upper/db/v4"
)

type upperdbSession interface {
	// ConnectionURL returns the DSN that was used to set up the adapter.
	ConnectionURL() upperdb.ConnectionURL

	// Name returns the name of the database.
	Name() string

	// Ping returns an error if the DBMS could not be reached.
	Ping() error

	// Collection receives a table name and returns a collection reference. The
	// information retrieved from a collection is cached.
	Collection(name string) upperdb.Collection

	// Collections returns a collection reference of all non system tables on the
	// database.
	Collections() ([]upperdb.Collection, error)

	// Save creates or updates a record.
	Save(record upperdb.Record) error

	// Get retrieves a record that matches the given condition.
	Get(record upperdb.Record, cond interface{}) error

	// Delete deletes a record.
	Delete(record upperdb.Record) error

	// Reset resets all the caching mechanisms the adapter is using.
	Reset()

	// Close terminates the currently active connection to the DBMS and clears
	// all caches.
	Close() error

	// Driver returns the underlying driver of the adapter as an interface.
	//
	// In order to actually use the driver, the `interface{}` value needs to be
	// casted into the appropriate types.
	//
	// Example:
	//  internalSQLDriver := sess.Driver().(*sql.DB)
	Driver() interface{}

	// Tx creates a transaction block on the default database context and passes
	// it to the function fn. If fn returns no error the transaction is commited,
	// else the transaction is rolled back. After being commited or rolled back
	// the transaction is closed automatically.
	Tx(fn func(sess upperdb.Session) error) error

	// TxContext creates a transaction block on the given context and passes it to
	// the function fn. If fn returns no error the transaction is commited, else
	// the transaction is rolled back. After being commited or rolled back the
	// transaction is closed automatically.
	TxContext(ctx context.Context, fn func(sess upperdb.Session) error, opts *sql.TxOptions) error

	// Context returns the context used as default for queries on this session
	// and for new transactions.  If no context has been set, a default
	// context.Background() is returned.
	Context() context.Context

	// WithContext returns a copy of the session that uses the given context as
	// default. Copies are safe to use concurrently but they're backed by the
	// same Session. You may close a copy at any point but that won't close the
	// parent session.
	WithContext(ctx context.Context) upperdb.Session
}

type dbSession struct {
	upperdbSession
	upperdb.SQL
}

func (sess *dbSession) Ctx(ctx context.Context) Session {
	return &dbSession{sess.WithContext(ctx), sess.SQL}
}

func (sess *dbSession) Transaction(ctx context.Context, fn func(sess Session) error) error {
	return sess.TxContext(ctx, func(sess upperdb.Session) error {
		return fn(&txSession{sess, sess.SQL()})
	}, nil)
}

func (sess *dbSession) InsertRecord(record *Record) error {
	return insertRecord(sess, record)
}
func (sess *dbSession) Commit() error {
	return nil
}
func (sess *dbSession) Rollback() error {
	return nil
}

type txSession struct {
	upperdbSession
	upperdb.SQL
}

func (sess *txSession) Ctx(ctx context.Context) Session {
	return &txSession{sess.WithContext(ctx), sess.SQL}
}

func (sess *txSession) Transaction(ctx context.Context, fn func(sess Session) error) error {
	return nil
	// return fn(&txSession{sess.WithContext(ctx)})
}

func (sess *txSession) InsertRecord(record *Record) error {
	return insertRecord(sess, record)
}

func NewDatabase(db upperdb.Session) Database {
	return &dbSession{db, db.SQL()}
}

type Session interface {
	upperdb.SQL
	upperdbSession
	InsertRecord(*Record) error
	Ctx(context.Context) Session
	Transaction(ctx context.Context, fn func(sess Session) error) error
}

type Database interface {
	upperdbSession
	Ctx(context.Context) Session
}
