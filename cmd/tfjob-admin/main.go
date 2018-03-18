package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"tfjob-admin/cmd/tfjob-admin/app"
	"tfjob-admin/cmd/tfjob-admin/app/options"
	//	"time"
)

func main() {
	var stop <-chan int

	s := options.NewServerRunOptions()
	s.AddFlags(pflag.CommandLine)

	pflag.Parse()

	if err := app.Run(s, stop); err != nil {
		fmt.Fprint(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
