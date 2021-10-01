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
	var wg sync.WaitGroup

	for i, node := range nodesList {
		wg.Add(1)

		go func(i int, node string) {
			defer wg.Done()

			msg := fmt.Sprintf("node=%d,time=%v,key=1", i, time.Now())
			cid, err := PutString(ctx, node, msg)
			if err != nil {
				log.Println(node, err)
				return
			}
			log.Println("PutString", node, msg)

			for j, lookupNode := range nodesList {
				if i == j {
					continue
				}
				startAt := time.Now()
				_, err = LookupString(ctx, lookupNode, cid)
				if err != nil {
					log.Println(lookupNode, err)
					continue
				}
				log.Printf("LookupString put=%v lookup=%v elapsed=%v i=%d j=%d\n", node, lookupNode, time.Since(startAt), i, j)
			}
		}(i, node)
	}

	wg.Wait()
}
