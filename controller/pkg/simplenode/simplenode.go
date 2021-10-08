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
