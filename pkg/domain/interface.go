package domain

import (
	"context"
	"user-generator/pkg/types"
)

type Provider interface {
	List(ctx context.Context) ([]*types.User, error)
}

type Syncer interface {
	Sync(context.Context, *types.User) error
}
