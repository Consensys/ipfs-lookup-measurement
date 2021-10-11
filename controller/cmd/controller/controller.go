package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/config"
	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/simplenode"
)

func main() {
	log.SetFlags(log.Lshortfile)

	ctx := context.Background()

	simpleNodesFile := flag.String("l", "nodes-list.out", "nodes list file")
	intervalSeconds := flag.Int("i", 0, "interval between each test")

	flag.Parse()
	nodesList := config.GetNodesList(*simpleNodesFile)

	for {
		simplenode.Experiment(ctx, nodesList)
		log.Println("one test is done")
		if *intervalSeconds == 0 {
			break
		}
		time.Sleep(time.Duration(*intervalSeconds) * time.Second)
	}
}
