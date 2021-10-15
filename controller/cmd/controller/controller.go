package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/config"
	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/server"
	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/simplenode"
)

func main() {
	simpleNodesFile := flag.String("l", "nodes-list.out", "nodes list file")
	intervalSeconds := flag.Int("i", 0, "interval between each test")

	flag.Parse()
	nodesList := config.GetNodesList(*simpleNodesFile)

	// Try to load key
	keyStr, err := ioutil.ReadFile(".key")
	if err != nil {
		fmt.Printf("error in getting the key: %v\n", err.Error())
		return
	}
	key, err := base64.RawStdEncoding.DecodeString(string(keyStr))
	if err != nil {
		fmt.Printf("error decoding key string: %v\n", err.Error())
		return
	}
	if len(key) != 32 {
		fmt.Printf("Wrong key size, expect 32, got: %v\n", len(key))
		return
	}
	// At start up, ask for list of node IDs.
	ids := make([]string, 0)
	for _, node := range nodesList {
		fmt.Printf("Start asking for node id from %v\n", node)
		id, err := server.RequestGetID(node, key)
		if err != nil {
			fmt.Printf("error getting node id for %v: %v\n", node, err.Error())
			return
		}
		fmt.Printf("Got node id for %v: %v\n", node, id)
		ids = append(ids, id)
	}
	// Ask every node to set IDs.
	for _, node := range nodesList {
		fmt.Printf("Start asking node %v to set up ids\n", node)
		out, err := server.RequestSetID(node, key, ids)
		if err != nil {
			fmt.Printf("error setting id for node %v: %v", node, err.Error())
			return
		}
		fmt.Printf("Got response for setting id for node %v: %v\n", node, out)
	}

	// Start the experiment.
	publish := 0
	max := len(nodesList) - 1
	for {
		simplenode.Experiment(publish, key, nodesList)
		log.Println("one test is done")
		if *intervalSeconds == 0 {
			break
		}
		publish++
		if publish > max {
			publish = 0
		}
		time.Sleep(time.Duration(*intervalSeconds) * time.Second)
	}
}
