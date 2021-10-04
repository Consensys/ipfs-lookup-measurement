// simplenode package will test with ipfs api
package simplenode

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Experiment publish message from one node, and lookup from all other nodes
func Experiment(ctx context.Context, nodesList []string) {
	var wg01 sync.WaitGroup
	for i, node := range nodesList {
		wg01.Add(1)

		go func(i int, node string) {
			defer wg01.Done()

			// publish string from nodes[i]
			msg := fmt.Sprintf("node=%d,time=%v,key=1", i, time.Now())
			cid, err := PutString(ctx, node, msg)
			if err != nil {
				log.Println(node, err)
				return
			}
			log.Println("PutString", node, msg)

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
					_, err = LookupString(ctx, lookupNode, cid)
					if err != nil {
						log.Println(lookupNode, err)
						return
					}
					log.Printf("LookupString put=%v lookup=%v elapsed=%v i=%d j=%d\n", node, lookupNode, time.Since(startAt), i, j)
				}(i, j, node, lookupNode)
			}
			wg02.Wait()

		}(i, node)
	}
	wg01.Wait()
}
