package pkg

import (
	"context"
	"user-export/pkg/api"
)

type UserProvider interface {
	List() ([]UserInterface, error)
}

type UserGenerator interface {
	Generate(context.Context, UserProvider) error
}

type UserInterface interface {
	ConvertCR() *api.User
}
