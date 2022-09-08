package ksgenerator

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"user-export/pkg"
	"user-export/pkg/api/v1alpha2"

	"k8s.io/client-go/tools/clientcmd"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ksGenerator struct {
	client rtclient.Client
}

func NewKSGenerator(options *Options) (pkg.UserGenerator, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", options.KubeConfigPath)
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
		err := ke.createUser(ctx, cr)
		if err != nil {
			return err
		}
		log.Println(fmt.Sprintf("create user %s success", cr.Name))
	}
	return nil
}
