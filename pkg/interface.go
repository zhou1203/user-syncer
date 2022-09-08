package pkg

import (
	"context"

	"user-export/pkg/api/v1alpha2"
)

type UserProvider interface {
	List() ([]User, error)
}

type UserGenerator interface {
	Generate(context.Context, UserProvider) error
}

type UserInterface interface {
	ConvertCR() *v1alpha2.User
}
