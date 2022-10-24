package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gojek/heimdall/v7/httpclient"
	"io/ioutil"
	"net/http"
	"time"
	"user-syncer/pkg/domain"
	"user-syncer/pkg/types"
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", h.Host, h.OrgPath), nil)
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
