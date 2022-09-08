package ksexport

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"user-export/pkg"
	"user-export/pkg/api"
)

type ksGenerator struct {
	client rest.Interface
}

func NewKubernetesExport(options *Options) (pkg.UserGenerator, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", options.KubeConfigPath)
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &ksGenerator{client: clientSet.RESTClient()}, nil
}

func (ke *ksGenerator) createUser(ctx context.Context, user *api.User) (api.User, error) {
	result := api.User{}
	err := ke.client.Post().
		Resource("users").
		Body(user).
		Do(ctx).
		Into(&result)
	return result, err
}

func (ke *ksGenerator) Generate(ctx context.Context, provider pkg.UserProvider) error {
	list, err := provider.List()
	if err != nil {
		return err
	}
	for _, u := range list {
		result, err := ke.createUser(ctx, u.ConvertCR())
		if err != nil {
			return err
		}
		log.Printf("created user: %s", result.GetName())
	}
	return nil
}
