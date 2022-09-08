package httpprovider

import "github.com/spf13/pflag"

type Options struct {
	Host    string `json:"host"`
	Path    string `json:"get_path"`
	Source  string `json:"source"`
	SysCode string `json:"sys_code"`
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("http provider", pflag.ContinueOnError)
	fs.StringVar(&o.Host, "host", o.Host, "http user provider host")
	fs.StringVar(&o.Path, "path", o.Path, "http user provider list path")
	fs.StringVar(&o.Source, "source", o.Source, "the user`s source")
	fs.StringVar(&o.SysCode, "syscode", o.SysCode, "http user provider sysCode")
	return fs
}
