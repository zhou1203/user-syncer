package syncer

import (
	"context"
	upperdb "github.com/upper/db/v4"
	"k8s.io/klog/v2"
	"reflect"
	"user-syncer/pkg/db"
	"user-syncer/pkg/domain"
	"user-syncer/pkg/types"
)

type orgDbSyncer struct {
	db db.Database
}

const (
	tableOrg = "org"

	columnOrgName     = "ORG_NAME"
	columnParentOrgID = "PARENT_ORG_ID"
)

func (ds *orgDbSyncer) Sync(ctx context.Context, obj interface{}) error {
	org := obj.(*types.Org)
	if status, err := ds.createOrUpdateInDB(ctx, org); err != nil {
		return err
	} else {
		switch status {
		case statusCreated:
			klog.Infof("Database: created org %s successful", org.ID)
		case statusUpdated:
			klog.Infof("Database: org existed, updated org %s successful", org.ID)
		case statusNoChange:
			klog.Infof("Database: org %s existed, no change", org.ID)
		}
	}
	return nil
}

func NewOrgDBSyncer(db db.Database) domain.Syncer {
	return &orgDbSyncer{db: db}
}

func (ds *orgDbSyncer) createOrUpdateInDB(ctx context.Context, org *types.Org) (status, error) {

	orgGet := &types.Org{}
	err := ds.db.Ctx(ctx).SelectFrom(tableOrg).Where(columnOrgID, org.ID).One(orgGet)
	if err != nil {
		if err == upperdb.ErrNoMoreRows {
			return statusCreated, ds.db.Ctx(ctx).InsertRecord(db.NewRecordWithObject(tableOrg, org))
		}
		return statusNoChange, err
	}

	if !reflect.DeepEqual(org, orgGet) {
		stmt := ds.db.Ctx(ctx).Update(columnOrgID)
		if org.OrgName != orgGet.OrgName {
			stmt = stmt.Set(columnOrgName, org.OrgName)
		}
		if org.ParentOrgID != orgGet.ParentOrgID {
			stmt = stmt.Set(columnParentOrgID, org.ParentOrgID)
		}

		_, err = stmt.Where(columnOrgID, org.ID).Exec()
		return statusUpdated, err
	}
	return statusNoChange, nil
}
