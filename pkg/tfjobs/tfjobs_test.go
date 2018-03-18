package tfjobs

import (
	"fmt"
)

func main() {
	tfjobclient, err := NewTFJobClient("", "./kube-proxy.kubeconfig")

	if err != nil {
		fmt.Print("测试没有通过")
		return
	}
	tfjoblist, err := tfjobclient.ListJobs("test", "default")

	if err != nil {
		fmt.Print("test no pass")
		return
	}
	fmt.Println(tfjoblist)

}
