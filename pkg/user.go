package pkg

import (
	"strconv"

	"user-generator/pkg/api/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type User struct {
	ID                int64  `json:"id"`
	Source            string `json:"-"`
	Name              string `json:"username"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"password"`
}

func (u User) ConvertCR() *v1alpha2.User {
	id := strconv.FormatInt(u.ID, 10)
	userCRD := &v1alpha2.User{
		TypeMeta: metav1.TypeMeta{
			Kind:       "User",
			APIVersion: "iam.kubesphere.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: u.Name,
			Labels: map[string]string{
				"iam.kubesphere.io/identify-provider": u.Source,
				"iam.kubesphere.io/origin-uid":        id,
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
