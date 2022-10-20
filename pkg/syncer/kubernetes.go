package syncer

import (
	"context"
	"log"
	"reflect"
	"user-generator/pkg/api/v1alpha2"
	"user-generator/pkg/domain"
	"user-generator/pkg/types"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8stypes "k8s.io/apimachinery/pkg/types"

	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ksSyncer struct {
	client rtclient.Client
}

func NewKSSyncer(client rtclient.Client) domain.Syncer {
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
		return true, ks.client.Update(ctx, user)
	}
	return true, nil
}

func (ks *ksSyncer) Sync(ctx context.Context, user *types.User) error {
	cr := ks.toObject(user)

	if user.Status == 0 {
		exist, err := ks.createOrUpdateUserInKS(ctx, cr)
		if err != nil {
			if errors.IsInternalError(err) {
				return err
			}
			log.Println(err)
		} else {
			if exist {
				log.Printf("user existed, update user %s success", cr.Name)
			} else {
				log.Printf("create user %s success", cr.Name)
			}
		}
	}
	return nil
}

func (ks *ksSyncer) toObject(u *types.User) *v1alpha2.User {
	userCRD := &v1alpha2.User{
		TypeMeta: metav1.TypeMeta{
			Kind:       "User",
			APIVersion: "iam.kubesphere.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: u.Name,
			Labels: map[string]string{
				"iam.kubesphere.io/identify-provider": u.Source,
				"iam.kubesphere.io/origin-uid":        u.ID,
			},
			Annotations: map[string]string{
				"ldap-manager/org-id": u.OrgID,
			},
		},
		Spec: v1alpha2.UserSpec{
			Email: u.Email,
		},
		Status: v1alpha2.UserStatus{
			State: "Active",
		},
	}
	return userCRD
}
