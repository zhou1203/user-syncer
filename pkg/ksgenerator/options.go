package ksgenerator

import (
	"github.com/spf13/pflag"
)

type Options struct {
	KubeConfigPath string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("k8s generator", pflag.ContinueOnError)
	flags.StringVar(&o.KubeConfigPath, "kubeconfig", o.KubeConfigPath, "kubeconfig file absolute path.")
	return flags
}
