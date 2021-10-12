package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/messaging"
)

func main() {
	log.SetFlags(log.Lshortfile)

	ipfsTestFolder := os.Getenv("PERFORMANCE_TEST_DIR")
	if ipfsTestFolder == "" {
		ipfsTestFolder = "/ipfs-tests"
	}

	err := os.Chdir(ipfsTestFolder)
	if err != nil {
		log.Fatalln(ipfsTestFolder, err)
	}

	cmd := flag.NewFlagSet("simple", flag.ExitOnError)
	portNumStr := cmd.String("p", "3030", "port number")

	log.Println("start listening at:", *portNumStr)

	http.HandleFunc("/", messaging.AgentHangler().ServeHTTP)
	log.Fatal(http.ListenAndServe(":"+*portNumStr, nil))
}
