package provider

import (
	"context"
	"user-generator/pkg"
	"user-generator/pkg/types"
)

type fakeProvider struct {
	options *Options
}

func NewFakeProvider(options *Options) pkg.Provider {
	return &fakeProvider{
		options: options,
	}
}

var users = []types.User{
	{
		ID:    "1",
		Name:  "fake-syncer-user-1",
		Email: "fakeprovider@kubesphere.io",
	},
	{
		ID:    "2",
		Name:  "fake-syncer-user-2",
		Email: "fakeprovider2@kubesphere.io",
	},
	{
		ID:    "3",
		Name:  "fake-syncer-user-3",
		Email: "fakeprovider3@kubesphere.io",
	},
	{
		ID:    "4",
		Name:  "fake-syncer-user-4",
		Email: "fakeprovider4@kubesphere.io",
	},
}

func (p *fakeProvider) List(ctx context.Context) ([]*types.User, error) {
	ui := make([]*types.User, 0)
	for _, v := range users {
		v.Source = p.options.Source
		ui = append(ui, &v)
	}
	return ui, nil
}
