package httpprovider

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/json"

	"user-generator/pkg"

	"github.com/gojek/heimdall/v7/httpclient"
)

type httpUserProvider struct {
	*Options
	httpClient *httpclient.Client
}

func NewHttpProvider(options *Options) (pkg.UserProvider, error) {
	provider := &httpUserProvider{
		Options: options,
	}

	provider.httpClient = httpclient.NewClient(httpclient.WithHTTPTimeout(30 * time.Second))

	return provider, nil

}

func (h *httpUserProvider) List() ([]*pkg.User, error) {
	users := make([]*pkg.User, 0)
	u, err := url.Parse(fmt.Sprintf("%s/%s", h.Host, h.Path))
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		URL:    u,
		Method: http.MethodGet,
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
