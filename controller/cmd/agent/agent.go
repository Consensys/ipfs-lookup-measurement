package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/messaging"

	api "github.com/ipfs/go-ipfs-api"
)

func main() {
	log.SetFlags(log.Lshortfile)

	ipfsTestFolder := os.Getenv("PERFORMANCE_TEST_DIR")
	if ipfsTestFolder == "" {
		ipfsTestFolder = "/ipfs-tests"
	}

	baseDir := os.Getenv("IPFS_PATH")
	fmt.Println("basedir is", baseDir)

	err := os.Chdir(ipfsTestFolder)
	if err != nil {
		log.Fatalln(ipfsTestFolder, err)
	}

	cmd := flag.NewFlagSet("simple", flag.ExitOnError)
	portNumStr := cmd.String("p", "3030", "port number")

	for {
		sh := api.NewLocalShell()
		if sh == nil {
			fmt.Println("error getting local shell")
		} else {
			fmt.Println("good at getting local shell")
			break
		}
		time.Sleep(1 * time.Second)
	}

	log.Println("start listening at:", *portNumStr)

	http.HandleFunc("/", messaging.AgentHangler().ServeHTTP)
	log.Fatal(http.ListenAndServe(":"+*portNumStr, nil))
}
