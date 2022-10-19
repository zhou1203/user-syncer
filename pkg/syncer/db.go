package syncer

import (
	"context"
	"user-generator/pkg"
	"user-generator/pkg/db"
)

type dbSyncer struct {
	db db.Database
}

func (d dbSyncer) Sync(ctx context.Context, provider pkg.Provider) error {
	//TODO implement me
	panic("implement me")
}

func NewDBSyncer(db db.Database) pkg.Syncer {
	return &dbSyncer{db: db}
}
