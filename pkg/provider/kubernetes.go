package provider

import (
	"context"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
	"user-generator/pkg"
	"user-generator/pkg/api/v1alpha2"
	"user-generator/pkg/types"
)

type ksProvider struct {
	client rtclient.Client
	source string
}

func (k *ksProvider) List(ctx context.Context) ([]*types.User, error) {
	users := make([]*types.User, 0)
	list := &v1alpha2.UserList{}
	err := k.client.List(ctx, list, nil)
	if err != nil {
		return nil, err
	}
	for _, user := range list.Items {
		var userObj types.User
		if user.Labels["iam.kubesphere.io/identify-provider"] == k.source {

			if user.Status.State == "Active" {
				userObj.Status = 1
			} else {
				userObj.Status = 0
			}
			userObj = types.User{
				ID:     user.Labels["iam.kubesphere.io/origin-uid"],
				Source: user.Labels["iam.kubesphere.io/identify-provider"],
				Name:   user.Name,
				Email:  user.Spec.Email,
			}
		}
		users = append(users, &userObj)
	}

	return users, nil
}

func NewKSProvider(client rtclient.Client) pkg.Provider {
	return &ksProvider{client: client}
}
