package domain

import (
	"context"
)

type Provider interface {
	List(ctx context.Context) ([]interface{}, error)
}

type Syncer interface {
	Sync(context.Context, interface{}) error
}
