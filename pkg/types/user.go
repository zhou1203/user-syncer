package types

import (
	"user-generator/pkg/api/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type User struct {
	ID      string `json:"id" db:"USER_ID"`
	LoginNo string `db:"LOGIN_NO"`
	Status  int    `json:"status" db:"APP_ACCT_STATUS"`
	Name    string `json:"username" db:"USER_NAME"`
	Email   string `json:"email" db:"EMAIL"`
	OrgID   string `json:"orgId" db:"ORG_ID"`
	Mobile  string `json:"mobile" db:"MOBILE"`
	Source  string `json:"-"`
}

func (u User) ConvertCR() *v1alpha2.User {
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
