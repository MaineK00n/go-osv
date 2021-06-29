package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MaineK00n/go-osv/cmd"
	"github.com/MaineK00n/go-osv/config"
)

func main() {
	var v = flag.Bool("v", false, "Show version")
	flag.Parse()
	if *v {
		fmt.Printf("go-osv-%s-%s\n", config.Version, config.Revision)
		os.Exit(0)
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
