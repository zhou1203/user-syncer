package provider

import "github.com/spf13/pflag"

type Options struct {
	Host     string
	UserPath string
	OrgPath  string
	Source   string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("provider http", pflag.ContinueOnError)
	fs.StringVar(&o.Host, "host", o.Host, "provider user syncer host")
	fs.StringVar(&o.UserPath, "user-path", o.UserPath, "provider user syncer list path")
	fs.StringVar(&o.OrgPath, "org-path", o.OrgPath, "provider org path")
	fs.StringVar(&o.Source, "source", o.Source, "the user`s source")
	return fs
}
