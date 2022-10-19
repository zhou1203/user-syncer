package pkg

import (
	"context"
	"user-generator/pkg/types"

	"user-generator/pkg/api/v1alpha2"
)

type Provider interface {
	List(ctx context.Context) ([]*types.User, error)
}

type Syncer interface {
	Sync(context.Context, Provider) error
}

type UserInterface interface {
	ConvertCR() *v1alpha2.User
}
