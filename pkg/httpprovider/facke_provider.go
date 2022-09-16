package httpprovider

import "user-generator/pkg"

type fakeProvider struct {
	options *Options
}

func NewFakeProvider(options *Options) pkg.UserProvider {
	return &fakeProvider{
		options: options,
	}
}

var users = []pkg.User{
	{
		ID:    1,
		Name:  "fake-provider-user-1",
		Email: "fakeprovider@kubesphere.io",
	},
	{
		ID:    2,
		Name:  "fake-provider-user-2",
		Email: "fakeprovider2@kubesphere.io",
	},
	{
		ID:    3,
		Name:  "fake-provider-user-3",
		Email: "fakeprovider3@kubesphere.io",
	},
	{
		ID:    4,
		Name:  "fake-provider-user-4",
		Email: "fakeprovider4@kubesphere.io",
	},
}

func (p *fakeProvider) List() ([]*pkg.User, error) {
	ui := make([]*pkg.User, 0)
	for _, v := range users {
		v.Source = p.options.Source
		ui = append(ui, &v)
	}
	return ui, nil
}
