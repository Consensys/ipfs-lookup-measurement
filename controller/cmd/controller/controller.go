package main

import (
	"context"
	"flag"
	"log"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/config"
	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/simplenode"
)

func main() {
	log.SetFlags(log.Lshortfile)

	ctx := context.Background()

	simpleNodesFile := flag.String("l", "nodes-list.out", "nodes list file")

	flag.Parse()
	nodesList := config.GetNodesList(*simpleNodesFile)
	simplenode.Experiment(ctx, nodesList)
}
