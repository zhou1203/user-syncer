package httpprovider

import "github.com/spf13/pflag"

type Options struct {
	Host    string `json:"host"`
	Path    string `json:"get_path"`
	SysCode string `json:"sys_code"`
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("http provider", pflag.ContinueOnError)
	fs.StringVar(&o.Host, "host", o.Host, "http host")
	fs.StringVar(&o.Path, "path", o.Path, "http list path")
	fs.StringVar(&o.SysCode, "syscode", o.SysCode, "http sysCode")
	return fs
}
