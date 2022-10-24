package provider

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"user-syncer/pkg/domain"
	"user-syncer/pkg/types"

	"github.com/gojek/heimdall/v7/httpclient"
)

type orgProvider struct {
	*Options
	httpClient *httpclient.Client
}

func NewOrgProvider(options *Options) (domain.Provider, error) {
	provider := &orgProvider{
		Options: options,
	}

	provider.httpClient = httpclient.NewClient(httpclient.WithHTTPTimeout(30 * time.Second))

	return provider, nil

}

func (h *orgProvider) List(ctx context.Context) ([]interface{}, error) {
	orgs := make([]*types.Org, 0)
	objs := make([]interface{}, 0)

	u := url.URL{
		Scheme: "http",
		Host:   h.Options.Host,
		Path:   h.Options.OrgPath,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &orgs)
	if err != nil {
		return nil, err
	}

	for _, v := range orgs {
		if v.Status == 0 {
			objs = append(objs, v)
		}
	}
	return objs, nil
}
