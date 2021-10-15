// simplenode package will test with ipfs api
package simplenode

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/server"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("controller")

// Experiment publish message from one node, and lookup from all other nodes
func Experiment(publish int, key []byte, nodesList []string) {
	// Publish string from node[publish]
	publisher := nodesList[publish]

	// Generate random content, 512 bytes.
	content := make([]byte, 512)
	rand.Read(content)

	// Request Publish
	cid, err := server.RequestPublish(publisher, key, content)
	if err != nil {
		log.Errorf("Error in publishing content from %v: %v", publisher, err.Error())
		return
	}
	log.Infof("Published content from %v with cid: %v", publisher, cid)
	// Need to wait till publish is finished.
	for i := 0; i < 60; i++ {
		time.Sleep(5 * time.Second)
		done, err := server.RequestCheck(publisher, key, cid)
		if err != nil {
			log.Errorf("Error in requesting a check to %v: %v", publisher, err.Error())
			return
		}
		if done {
			log.Infof("Publish from %v is done", publisher)
			break
		}
		if i == 59 {
			log.Errorf("Error in publishing from %v", publisher)
			return
		}
		log.Infof("Publish from %v in progress...", publisher)
	}
	// Start lookup from every other node.
	var wg sync.WaitGroup
	for i, lookupNode := range nodesList {
		if i == publish {
			continue
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, lookupNode string) {
			defer wg.Done()
			// First do a disconnection to avoid using bitswap
			out, err := server.RequestDisconnect(lookupNode, key)
			if err != nil {
				log.Errorf("Error requesting disconnection to %v: %v", lookupNode, err.Error())
				return
			}
			log.Infof("Response of disconnection from %v is: %v", lookupNode, out)
			err = server.RequestLookup(lookupNode, key, cid)
			if err != nil {
				log.Errorf("Error requesting lookup to %v: %v", lookupNode, err.Error())
				return
			}
			log.Infof("Start lookup %v from %v", cid, lookupNode)
			// Need to wait till lookup is finished.
			for i := 0; i < 30; i++ {
				time.Sleep(5 * time.Second)
				done, err := server.RequestCheck(lookupNode, key, cid)
				if err != nil {
					log.Errorf("Error in requesting a check to %v: %v", lookupNode, err.Error())
					return
				}
				if done {
					log.Infof("Lookup from %v is done", lookupNode)
					break
				}
				if i == 29 {
					log.Errorf("Error in lookup from %v", lookupNode)
					return
				}
				log.Infof("Lookup from %v in progress...", lookupNode)
			}
		}(&wg, lookupNode)
	}
	// Wait till all lookup is finished.
	wg.Wait()
	// Clean
	for _, node := range nodesList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, node string) {
			defer wg.Done()
			out, err := server.RequestClean(node, key, cid)
			if err != nil {
				log.Errorf("Error in requesting clean of %v to %v: %v", cid, node, err.Error())
			} else {
				log.Infof("Response of clean of %v from %v is: %v", cid, node, out)
			}
		}(&wg, node)
	}
	wg.Wait()
	log.Infof("clean is done")
}
