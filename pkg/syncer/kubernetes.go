package syncer

import (
	"context"
	"k8s.io/klog/v2"
	"reflect"
	"user-syncer/pkg/api/v1alpha2"
	"user-syncer/pkg/domain"
	"user-syncer/pkg/types"

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

func (ks *ksSyncer) createOrUpdateUserInKS(ctx context.Context, user *v1alpha2.User) (status, error) {
	userGet := &v1alpha2.User{}
	err := ks.client.Get(ctx, k8stypes.NamespacedName{Name: user.Name}, userGet)
	if err != nil {
		if errors.IsNotFound(err) {
			return statusCreated, ks.client.Create(ctx, user)
		}
		return statusNoChange, err
	}

	if userGet.Annotations["ldap-manager/org-id"] != user.Annotations["ldap-manager/org-id"] ||
		userGet.Labels["iam.kubesphere.io/identify-provider"] != user.Labels["iam.kubesphere.io/identify-provider"] ||
		userGet.Labels["iam.kubesphere.io/origin-uid"] != user.Labels["iam.kubesphere.io/origin-uid"] ||
		!reflect.DeepEqual(userGet.Spec, user.Spec) ||
		!reflect.DeepEqual(userGet.Status, user.Status) {
		return statusUpdated, ks.client.Update(ctx, user)
	}
	return statusNoChange, nil
}

func (ks *ksSyncer) Sync(ctx context.Context, obj interface{}) error {
	user := obj.(*types.User)
	cr := ks.toObject(user)

	if user.Status == 0 {
		status, err := ks.createOrUpdateUserInKS(ctx, cr)
		if err != nil {
			return err
		} else {
			switch status {
			case statusCreated:
				klog.Infof("Kubernetes: created user %s successful", user.Name)
			case statusUpdated:
				klog.Infof("Kubernetes: user existed, updated user %s successful", user.Name)
			case statusNoChange:
				klog.Infof("Kubernetes: user existed, user %s no change", user.Name)
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
