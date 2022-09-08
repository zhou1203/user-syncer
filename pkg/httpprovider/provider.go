package httpprovider

import (
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"net/url"
	"time"
	"user-export/pkg"

	"github.com/gojek/heimdall/v7/httpclient"
)

const source = "ldap"

type httpUserProvider struct {
	*Options
	httpClient *httpclient.Client
}

func NewHttpProvider(options *Options) (pkg.UserProvider, error) {
	provider := &httpUserProvider{
		Options: options,
	}

	provider.httpClient = httpclient.NewClient(httpclient.WithHTTPTimeout(1000 * time.Millisecond))

	return provider, nil

}

func (h *httpUserProvider) List() ([]pkg.UserInterface, error) {
	users := make([]user, 0)
	userInterfaces := make([]pkg.UserInterface, 0)
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
		u.Source = source
		userInterfaces = append(userInterfaces, &u)
	}

	return userInterfaces, nil

}
