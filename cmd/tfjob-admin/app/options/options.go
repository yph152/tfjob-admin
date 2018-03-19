package options

import (
	"github.com/spf13/pflag"
	//	"fmt"
)

type ServerRunOptions struct {
	Port       int
	Kubeconfig string
}

func NewServerRunOptions() *ServerRunOptions {
	opt := &ServerRunOptions{}

	return opt
}

func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&s.Port, "port", 8899, "specify port number - 8899 by  default")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", "/root/.kube/kubeconfig", "kubeconfig for k8s")
}