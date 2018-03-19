package api

import (
	"fmt"
	"github.com/caicloud/nirvana"
	v1 "tfjob-admin/pkg/api/v1/descriptor"
	v1handlers "tfjob-admin/pkg/api/v1/handlers"
)

func Initialize(s *nirvana.Config, kubeconfig, storagehost string) {
	v1.Initialize(s)
	v1handlers.NewClient(kubeconfig, storagehost)
	fmt.Println("vim-go")
}
