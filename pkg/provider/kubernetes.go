package provider

import (
	"context"
	"strconv"
	"user-syncer/pkg/api/v1alpha2"
	"user-syncer/pkg/domain"
	"user-syncer/pkg/types"

	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ksProvider struct {
	client rtclient.Client
	source string
}

func (k *ksProvider) List(ctx context.Context) ([]interface{}, error) {
	users := make([]interface{}, 0)
	var list v1alpha2.UserList
	err := k.client.List(ctx, &list, rtclient.MatchingLabels{"iam.kubesphere.io/identify-provider": k.source})
	if err != nil {
		return nil, err
	}
	for _, user := range list.Items {
		var userObj types.User
		userID, err := strconv.ParseInt(user.Labels["iam.kubesphere.io/origin-uid"], 10, 64)
		if err != nil {
			return nil, err
		}
		userObj = types.User{
			ID:      userID,
			LoginNo: user.Name,
			OrgID:   user.Annotations["ldap-manager/org-id"],
			Name:    user.Name,
			Email:   user.Spec.Email,
		}
		if user.Status.State == "Active" {
			userObj.Status = 1
		} else {
			userObj.Status = 0
		}

		users = append(users, &userObj)
	}

	return users, nil
}

func NewKSProvider(client rtclient.Client, source string) domain.Provider {
	return &ksProvider{client: client, source: source}
}
