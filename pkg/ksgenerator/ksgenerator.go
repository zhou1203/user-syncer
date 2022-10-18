package ksgenerator

import (
	"context"
	"fmt"
	"log"
	"user-generator/pkg"
	"user-generator/pkg/api/v1alpha2"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ksGenerator struct {
	client rtclient.Client
}

func NewKSGenerator(options *Options) (pkg.UserGenerator, error) {
	var restConfig *rest.Config
	var err error
	if options.KubeConfigPath == "InCluster" {
		restConfig, err = rest.InClusterConfig()
	} else {
		restConfig, err = clientcmd.BuildConfigFromFlags("", options.KubeConfigPath)
	}
	if err != nil {
		return nil, err
	}

	sch := scheme.Scheme
	err = v1alpha2.AddToScheme(sch)
	if err != nil {
		return nil, err
	}
	client, err := rtclient.New(restConfig, rtclient.Options{
		Scheme: sch,
	})
	if err != nil {
		return nil, err
	}

	return &ksGenerator{client: client}, nil
}

func (ke *ksGenerator) createUser(ctx context.Context, user *v1alpha2.User) error {
	return ke.client.Create(ctx, user)
}

func (ke *ksGenerator) Generate(ctx context.Context, provider pkg.UserProvider) error {
	list, err := provider.List()
	if err != nil {
		return err
	}
	for _, u := range list {
		cr := u.ConvertCR()
		if u.Status == 0 {
			err := ke.createUser(ctx, cr)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(fmt.Sprintf("create user %s success", cr.Name))
			}
		}
	}
	return nil
}
