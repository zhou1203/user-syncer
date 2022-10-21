package syncer

import (
	"context"
	"k8s.io/klog/v2"
	"reflect"
	"user-syncer/pkg/db"
	"user-syncer/pkg/domain"
	"user-syncer/pkg/types"

	upperdb "github.com/upper/db/v4"
)

type dbSyncer struct {
	db db.Database
}

type status string

const (
	tableUser = "user"

	statusCreated  status = "created"
	statusUpdated  status = "updated"
	statusNoChange status = "no change"

	columnUserID        = "USER_ID"
	columnUserName      = "USER_NAME"
	columnLoginNo       = "LOGIN_NO"
	columnOrgID         = "ORG_ID"
	columnEmail         = "EMAIL"
	columnMobile        = "MOBILE"
	columnAppAcctStatus = "APP_ACCT_STATUS"
)

func (ds *dbSyncer) Sync(ctx context.Context, user *types.User) error {
	if status, err := ds.createOrUpdateInDB(ctx, user); err != nil {
		return err
	} else {
		switch status {
		case statusCreated:
			klog.Infof("Database: created user %s successful", user.Name)
		case statusUpdated:
			klog.Infof("Database: user existed, updated user %s successful", user.Name)
		case statusNoChange:
			klog.Infof("Database: user existed, user %s no change", user.Name)
		}
	}
	return nil
}

func NewDBSyncer(db db.Database) domain.Syncer {
	return &dbSyncer{db: db}
}

func (ds *dbSyncer) createOrUpdateInDB(ctx context.Context, user *types.User) (status, error) {
	if user.Status == 0 {
		user.Status = 1
		user.LoginNo = user.Name
	} else {
		return statusNoChange, nil
	}

	userGet := &types.User{}
	err := ds.db.Ctx(ctx).SelectFrom(tableUser).Where("USER_ID", user.ID).One(userGet)
	if err != nil {
		if err == upperdb.ErrNoMoreRows {
			return statusCreated, ds.db.Ctx(ctx).InsertRecord(db.NewRecordWithObject(tableUser, user))
		}
		return statusNoChange, err
	}

	user.Source = userGet.Source

	if !reflect.DeepEqual(user, userGet) {
		stmt := ds.db.Ctx(ctx).Update(tableUser)
		if user.Name != userGet.Name {
			stmt = stmt.Set(columnUserName, user.Name)
		}
		if user.LoginNo != userGet.LoginNo {
			stmt = stmt.Set(columnLoginNo, user.LoginNo)
		}
		if user.OrgID != userGet.OrgID {
			stmt = stmt.Set(columnOrgID, user.OrgID)
		}
		if user.Email != userGet.Email {
			stmt = stmt.Set(columnEmail, user.Email)
		}
		if user.Mobile != userGet.Mobile {
			stmt = stmt.Set(columnMobile, user.Mobile)
		}
		if user.Status != userGet.Status {
			stmt = stmt.Set(columnAppAcctStatus, user.Status)
		}

		_, err = stmt.Where(columnUserID, user.ID).Exec()
		return statusUpdated, err
	}
	return statusNoChange, nil
}
