package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/server"
)

func main() {
	ipfsTestFolder := os.Getenv("PERFORMANCE_TEST_DIR")
	if ipfsTestFolder == "" {
		ipfsTestFolder = "/ipfs-tests"
	}

	err := os.Chdir(ipfsTestFolder)
	if err != nil {
		fmt.Printf("error in chdir: %v\n", err.Error())
		return
	}

	cmd := flag.NewFlagSet("simple", flag.ExitOnError)
	portNumStr := cmd.String("p", "3030", "port number")
	key, err := ioutil.ReadFile(".key")
	if err != nil {
		fmt.Printf("error in getting the key: %v\n", err.Error())
		return
	}
	server.NewServer(context.Background(), ":"+*portNumStr, string(key))
}
