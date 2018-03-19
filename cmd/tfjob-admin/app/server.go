package app

import (
	"fmt"
	"github.com/caicloud/nirvana"
	"github.com/caicloud/nirvana/plugins/metrics"
	"github.com/caicloud/nirvana/plugins/profiling"
	"tfjob-admin/cmd/tfjob-admin/app/options"
	"tfjob-admin/pkg/api"
	"time"
)

const (
	defaultPort               = 8888
	DefaultTerminationTimeout = 5 * time.Minute
)

func Run(opt *options.ServerRunOptions, stopCh <-chan int) error {
	config := nirvana.NewDefaultConfig()

	config.Configure(
		metrics.Path("/metrics"),
		profiling.Path("/debug/pprof"),
		profiling.Contention(true),
		nirvana.Port(defaultPort),
	)

	api.Initialize(config, opt.Kubeconfig, opt.StorageHost)

	if err := nirvana.NewServer(config).Serve(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
