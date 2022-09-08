package httpprovider

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"user-export/pkg/api"
)

type user struct {
	ID     string `json:"id"`
	Source string `json:"-"`
	Name   string `json:"username"`
	Email  string `json:"email"`
}

func (u *user) ConvertCR() *api.User {
	userCRD := &api.User{
		TypeMeta: metav1.TypeMeta{
			Kind:       "users",
			APIVersion: "iam.kubesphere.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: u.Name,
			Annotations: map[string]string{
				"kubesphere.io/creator": "admin",
				// TODO add provider relevant ID, and provider name
			},
			Finalizers: []string{
				"finalizers.kubesphere.io/users",
			},
		},
		Spec: api.UserSpec{
			Email: u.Email,
			Lang:  "zh",
		},
		Status: api.UserStatus{
			State: "Active",
		},
	}
	return userCRD
}
