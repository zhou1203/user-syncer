package provider

import (
	"context"

	"user-syncer/pkg/domain"
	"user-syncer/pkg/types"
)

type fakeProvider struct {
	options *Options
}

type fakeOrgProvider struct {
}

func (f fakeOrgProvider) List(ctx context.Context) ([]interface{}, error) {
	objs := make([]interface{}, 0)
	for _, o := range orgs {
		objs = append(objs, o)
	}
	return objs, nil
}

func NewFakeOrgProvider() domain.Provider {
	return &fakeOrgProvider{}
}

func NewFakeProvider(options *Options) domain.Provider {
	return &fakeProvider{
		options: options,
	}
}

var orgs = []*types.Org{
	{
		ID:          "51111",
		OrgName:     "test-org-1",
		ParentOrgID: "51112",
	},
	{
		ID:          "51112",
		OrgName:     "test-org-2",
		ParentOrgID: "51113",
	},
}

var users = []*types.User{
	{
		ID:     1,
		Name:   "fake-syncer-user-1",
		Email:  "fakeprovider@kubesphere.io",
		OrgID:  "51111",
		Status: 0,
	},
	{
		ID:     2,
		Name:   "fake-syncer-user-2",
		Email:  "fakeprovider2@kubesphere.io",
		OrgID:  "51111",
		Status: 0,
	},
	{
		ID:     3,
		Name:   "fake-syncer-user-3",
		Email:  "fakeprovider3@kubesphere.io",
		OrgID:  "51113",
		Status: 0,
	},
	{
		ID:     4,
		Name:   "fake-syncer-user-4",
		Email:  "fakeprovider4@kubesphere.io",
		OrgID:  "51111",
		Status: 0,
	},
	{
		ID:     5,
		Name:   "fake-syncer-user-5",
		Email:  "fakeprovider5@kubesphere.io",
		OrgID:  "51112",
		Status: 0,
	},
	{
		ID:     6,
		Name:   "fake-syncer-user-6",
		Email:  "fakeprovider6@kubesphere.io",
		OrgID:  "51112",
		Status: 0,
	},
}

func (p *fakeProvider) List(ctx context.Context) ([]interface{}, error) {
	ui := make([]interface{}, 0)
	for _, v := range users {
		v.Source = p.options.Source
		ui = append(ui, v)
	}
	return ui, nil
}
