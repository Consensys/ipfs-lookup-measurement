package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/config"
	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/simplenode"
)

func main() {
	log.SetFlags(log.Lshortfile)

	parseCommandLine()
}

func parseCommandLine() {

	simpleCmd := flag.NewFlagSet("simple", flag.ExitOnError)
	simpleNodesFile := simpleCmd.String("l", "nodes-list.out", "nodes list file")

	if len(os.Args) < 2 {
		subCommandUsage()
	}

	switch os.Args[1] {
	case "simple":
		simpleCmd.Parse(os.Args[2:])
		nodesList := config.GetNodesList(*simpleNodesFile)
		simplenode.Experiment(context.Background(), nodesList)
	default:
		subCommandUsage()
	}

}

func subCommandUsage() {
	fmt.Println("expected 'simple' or 'dht' subcommands")
	os.Exit(1)
}
