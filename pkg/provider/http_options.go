package provider

import "github.com/spf13/pflag"

type Options struct {
	Host   string `json:"host"`
	Path   string `json:"get_path"`
	Source string `json:"source"`
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("provider http", pflag.ContinueOnError)
	fs.StringVar(&o.Host, "host", o.Host, "provider user syncer host")
	fs.StringVar(&o.Path, "path", o.Path, "provider user syncer list path")
	fs.StringVar(&o.Source, "source", o.Source, "the user`s source")
	return fs
}
