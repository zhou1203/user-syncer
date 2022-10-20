package provider

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"user-generator/pkg/domain"
	"user-generator/pkg/types"

	"k8s.io/klog/v2"

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

func (h *httpUserProvider) List(ctx context.Context) ([]*types.User, error) {
	users := make([]*types.User, 0)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", h.Host, h.Path), nil)
	if err != nil {
		klog.Error(err)
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
		if strings.Contains(u.Name, "_") {
			u.Name = strings.Replace(u.Name, "_", "-", -1)
		}
		u.Source = h.Source
	}

	return users, nil
}
