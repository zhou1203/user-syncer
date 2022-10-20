package domain

import (
	"user-generator/pkg/api/v1alpha2"

	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func NewKubernetesClient(options *KubeOptions) (rtclient.Client, error) {
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
	return client, nil
}

func (o *KubeOptions) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("k8s syncer", pflag.ContinueOnError)
	flags.StringVar(&o.KubeConfigPath, "kubeconfig", o.KubeConfigPath, "kubeconfig file absolute path.")
	return flags
}

type KubeOptions struct {
	KubeConfigPath string
}

func NewKubeOptions() *KubeOptions {
	return &KubeOptions{}
}
