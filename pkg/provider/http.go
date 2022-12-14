package provider

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"user-syncer/pkg/domain"
	"user-syncer/pkg/types"

	"k8s.io/apimachinery/pkg/util/json"

	"github.com/gojek/heimdall/v7/httpclient"
)

type httpUserProvider struct {
	*Options
	httpClient *httpclient.Client
}

func NewHttpProvider(options *Options) (domain.Provider, error) {
	provider := &httpUserProvider{
		Options: options,
	}

	provider.httpClient = httpclient.NewClient(httpclient.WithHTTPTimeout(30 * time.Second))

	return provider, nil

}

func (h *httpUserProvider) List(ctx context.Context) ([]interface{}, error) {
	objs := make([]interface{}, 0)
	users := make([]*types.User, 0)

	u := url.URL{
		Scheme: "http",
		Host:   h.Options.Host,
		Path:   h.Options.UserPath,
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
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Status == 0 {
			u.Name = strings.Replace(u.Name, "_", "-", -1)
			u.Source = h.Source
			objs = append(objs, u)
		}
	}

	return objs, nil
}
