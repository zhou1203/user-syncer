package syncer

import (
	"context"
	"fmt"
	upperdb "github.com/upper/db/v4"
	"k8s.io/apimachinery/pkg/api/errors"
	"log"
	"reflect"
	"user-generator/pkg"
	"user-generator/pkg/api/v1alpha2"
	"user-generator/pkg/db"
	"user-generator/pkg/types"

	k8stypes "k8s.io/apimachinery/pkg/types"

	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ksSyncer struct {
	client   rtclient.Client
	dbClient db.Database
}

const (
	tableUser = "user"
)

func NewKSSyncer(client rtclient.Client) pkg.Syncer {
	return &ksSyncer{client: client}
}

func (ks *ksSyncer) createOrUpdateUserInKS(ctx context.Context, user *v1alpha2.User) (bool, error) {
	userGet := &v1alpha2.User{}
	err := ks.client.Get(ctx, k8stypes.NamespacedName{Name: user.Name}, userGet)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, ks.client.Create(ctx, user)
		}
		return false, err
	}

	if !reflect.DeepEqual(user, userGet) {
		return false, ks.client.Update(ctx, user)
	}
	return true, nil
}

func (ks *ksSyncer) createOrUpdateInDB(ctx context.Context, user *types.User) (bool, error) {
	if user.Status == 0 {
		user.Status = 1
	}

	userGet := &types.User{}
	err := ks.dbClient.Ctx(ctx).SelectFrom(tableUser).Where("USER_ID", user.ID).One(userGet)
	if err != nil {
		if err == upperdb.ErrNoMoreRows {

			return false, ks.dbClient.Ctx(ctx).InsertRecord(db.NewRecordWithObject(tableUser, user))
		}
		return false, err
	}
	stmt := ks.dbClient.Ctx(ctx).Update(tableUser)
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

func (ks *ksSyncer) Sync(ctx context.Context, provider pkg.Provider) error {
	list, err := provider.List(ctx)
	if err != nil {
		return err
	}
	for _, u := range list {
		cr := u.ConvertCR()
		if u.Status == 0 {
			exist, err := ks.createOrUpdateUserInKS(ctx, cr)
			if err != nil {
				if errors.IsInternalError(err) {
					return err
				}
				log.Println(err)
			} else {
				if exist {
					log.Println(fmt.Sprintf("user existed, update user %s success", cr.Name))
				} else {
					log.Println(fmt.Sprintf("create user %s success", cr.Name))
				}
			}
			if exist, err := ks.createOrUpdateInDB(ctx, u); err != nil {
				log.Panicln(err)
				return err
			} else {
				if exist {
					log.Println(fmt.Sprintf("database: user existed, updated user %s success", u.Name))
				} else {
					log.Println(fmt.Sprintf("database: createed user %s success", u.Name))
				}
			}
		}
	}

	return nil
}
