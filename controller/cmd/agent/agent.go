package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/messaging"
)

func main() {
	log.SetFlags(log.Lshortfile)

	cmd := flag.NewFlagSet("simple", flag.ExitOnError)
	portNumStr := cmd.String("p", "3030", "port number")

	log.Println("start listening at:", *portNumStr)

	http.HandleFunc("/", messaging.AgentHangler().ServeHTTP)
	log.Fatal(http.ListenAndServe(":"+*portNumStr, nil))
}
