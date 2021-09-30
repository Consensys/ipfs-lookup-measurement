// simplenode package will test with ipfs api
package simplenode

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Experiment publish message from one node, and lookup from all other nodes
func Experiment(ctx context.Context, nodesList []string) {
	for i, node := range nodesList {
		msg := fmt.Sprintf("node=%d,time=%v,key=1", i, time.Now())
		cid, err := PutString(ctx, node, msg)
		if err != nil {
			log.Println(node, err)
			continue
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
			log.Println("LookupString", lookupNode, time.Since(startAt))
		}
	}
}
