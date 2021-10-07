// simplenode package will test with ipfs api
package simplenode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/messaging"
)

// Experiment publish message from one node, and lookup from all other nodes
func Experiment(ctx context.Context, nodesList []string) {
	var wg01 sync.WaitGroup
	for i, node := range nodesList {
		wg01.Add(1)

		go func(i int, node string) {
			defer wg01.Done()

			// publish string from nodes[i]
			m := messaging.RequestMessage{}
			m.IntOption1 = 3
			m.StrOption1 = fmt.Sprintf("node=%d,time=%v,key=1", i, time.Now())

			err := postCall("publish", node, m)
			if err != nil {
				log.Println(node, err)
				return
			}
			log.Println("publish", node)

			// lookup on all nodes except nodes[i]
			var wg02 sync.WaitGroup
			for j, lookupNode := range nodesList {
				wg02.Add(1)

				go func(i int, j int, node string, lookupNode string) {
					defer wg02.Done()

					if i == j {
						return
					}
					startAt := time.Now()
					err = postCall("lookup", node, m)
					if err != nil {
						log.Println(lookupNode, err)
						return
					}
					log.Printf("lookup put=%v lookup=%v elapsed=%v i=%d j=%d\n", node, lookupNode, time.Since(startAt), i, j)
				}(i, j, node, lookupNode)
			}
			wg02.Wait()

		}(i, node)
	}
	wg01.Wait()
}

func postCall(cmd string, node string, m messaging.RequestMessage) error {
	m.Cmd = cmd
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = http.Post(node, "application/json", bytes.NewBuffer(b))
	return err
}
