// simplenode package will test with ipfs api
package simplenode

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ConsenSys/ipfs-lookup-measurement/controller/pkg/messaging"
)

// Experiment publish message from one node, and lookup from all other nodes
func Experiment(ctx context.Context, publish int, nodesList []string) {
	// Publish string from node[publish]
	node := nodesList[publish]
	m := messaging.RequestMessage{}
	m.IntOption1 = 3
	m.StrOption1 = fmt.Sprintf("node=%d,time=%v,key=1", publish, time.Now())

	err := postCall("publish", node, m)
	if err != nil {
		log.Println(node, err)
		return
	}
	log.Println("publish", node)
	// Need to wait till publish is finished.
	for i := 0; i < 30; i++ {
		time.Sleep(5 * time.Second)
		err = postCall("check", node, m)
		if err == nil {
			log.Println("publish is done")
			break
		}
		if i == 29 {
			log.Println("error in publishing.")
			return
		}
	}
	// Start lookup from every other node.
	var wg sync.WaitGroup
	for i, lookupNode := range nodesList {
		if i == publish {
			return
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, lookupNode string) {
			defer wg.Done()
			err = postCall("lookup", lookupNode, m)
			if err != nil {
				log.Println(lookupNode, err)
				return
			}
			// Need to wait till lookup is finished.
			for i := 0; i < 30; i++ {
				time.Sleep(5 * time.Second)
				err = postCall("check", lookupNode, m)
				if err == nil {
					log.Println("lookup is done")
					break
				}
				if i == 29 {
					log.Println("error in lookup.")
					return
				}
			}
			log.Printf("lookup put=%v lookup=%v is done\n", node, lookupNode)
		}(&wg, lookupNode)
	}
	// Wait till all lookup is finished.
	wg.Wait()
	// Clean
	for _, lookupNode := range nodesList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, lookupNode string) {
			defer wg.Done()
			err = postCall("swarmdisconnect", lookupNode, m)
			if err != nil {
				log.Println("swarmdisconnect", lookupNode, err)
			}
		}(&wg, lookupNode)
	}
	wg.Wait()
	log.Println("clean is done")
}

func postCall(cmd string, node string, m messaging.RequestMessage) error {
	m.Cmd = cmd
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	buf, err := encryptStream(b)
	if err != nil {
		return err
	}

	_, err = http.Post(node, "application/json", buf)
	return err
}

func encryptStream(bs []byte) (buf io.Reader, err error) {
	key, err := ioutil.ReadFile(".key")
	if err != nil {
		return
	}
	key = key[:32]

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.

	return bytes.NewBuffer(gcm.Seal(nonce, nonce, bs, nil)), nil
}
