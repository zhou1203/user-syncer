package syncer

import (
	"context"
	"log"
	"user-syncer/pkg/db"
	"user-syncer/pkg/domain"
	"user-syncer/pkg/types"

	upperdb "github.com/upper/db/v4"
)

type dbSyncer struct {
	db db.Database
}

const (
	tableUser = "user"
)

func (ds *dbSyncer) Sync(ctx context.Context, user *types.User) error {
	if exist, err := ds.createOrUpdateInDB(ctx, user); err != nil {
		log.Panicln(err)
		return err
	} else {
		if exist {
			log.Printf("database: user existed, updated user %s success", user.Name)
		} else {
			log.Printf("database: createed user %s success", user.Name)
		}
	}
	return nil
}

func NewDBSyncer(db db.Database) domain.Syncer {
	return &dbSyncer{db: db}
}

func (ds *dbSyncer) createOrUpdateInDB(ctx context.Context, user *types.User) (bool, error) {
	if user.Status == 0 {
		user.Status = 1
	}

	userGet := &types.User{}
	err := ds.db.Ctx(ctx).SelectFrom(tableUser).Where("USER_ID", user.ID).One(userGet)
	if err != nil {
		if err == upperdb.ErrNoMoreRows {
			return false, ds.db.Ctx(ctx).InsertRecord(db.NewRecordWithObject(tableUser, user))
		}
		return false, err
	}
	stmt := ds.db.Ctx(ctx).Update(tableUser)
	if user.Name != userGet.Name {
		stmt = stmt.Set("USER_NAME", user.Name)
	}
	if user.LoginNo != userGet.LoginNo {
		stmt = stmt.Set("LOGON_NO", user.LoginNo)
	}
	if user.OrgID != userGet.OrgID {
		stmt = stmt.Set("ORG_ID", user.OrgID)
	}
	if user.Email != userGet.Email {
		stmt = stmt.Set("EMAIL", user.Email)
	}
	if user.Mobile != userGet.Mobile {
		stmt = stmt.Set("MOBILE", user.Mobile)
	}
	if user.Status != userGet.Status {
		stmt = stmt.Set("APP_ACCT_STATUS", user.Status)
	}

	_, err = stmt.Where("USER_ID", user.ID).Exec()
	return true, err
}
