package messaging

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	api "github.com/ipfs/go-ipfs-api"
)

type RequestMessage struct {
	Cmd        string
	IntOption1 int
	StrOption1 string
}

var syncMap sync.Map

func AgentHangler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// apply restrictions: post only
			if r.ContentLength == 0 {
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintln(w, "failed")
				return
			}

			// decrypt r.body

			// parsing json
			var m RequestMessage
			// err := json.NewDecoder(encryptStream(r.Body)).Decode(&m)
			plaintext, err := decryptStream(r.Body)
			if err != nil {
				log.Println(err)
				return
			}
			err = json.Unmarshal(plaintext, &m)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "failed")
				log.Println(err)
				return
			}

			if m.Cmd == "check" {
				if check(m) == nil {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusFound)
				}
				return
			}

			// handling request task
			err = taskHandler(m)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusExpectationFailed)
				fmt.Fprintln(w, "failed")
				return
			}

			// response
			// w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// m = RequestMessage{"cmd1", 1, "s2"}
			// json.NewEncoder(w).Encode(m)

			fmt.Fprintln(w, "ok")
		})
}

func decryptStream(r io.Reader) (plaintext []byte, err error) {
	key, err := ioutil.ReadFile(".key")
	if err != nil {
		return
	}
	key = key[:32]

	// if our program was unable to read the file
	// print out the reason why it can't
	if err != nil {
		return
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return
	}

	nonceSize := gcm.NonceSize()

	ciphertext, err := ioutil.ReadAll(r)
	if len(ciphertext) < nonceSize {
		return
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	return
}

func taskHandler(m RequestMessage) error {
	switch m.Cmd {
	case "publish":
		return publish(m)
	case "lookup":
		return lookup(m)
	case "swarmdisconnect":
		return swarmDisconnect(m)
	}
	return errors.New("command is invalid")
}

func swarmDisconnect(m RequestMessage) error {
	ipfs := getIPFS()

	for _, node := range getPeers(ipfs) {
		if _, ok := syncMap.Load(node); ok || node == "" {
			continue
		}
		cmdLine := fmt.Sprintf("%s swarm disconnect %s", ipfs, node)
		out, _ := exec.Command("sh", "-xc", cmdLine).CombinedOutput()
		log.Print(string(out))
	}

	return nil
}

func getIPFS() string {
	ipfs := os.Getenv("IPFS")
	if ipfs == "" {
		ipfs = "/app/go-ipfs/cmd/ipfs/ipfs"
	}
	return ipfs
}

func getPeers(ipfs string) (rtn []string) {
	cmdLine := fmt.Sprintf("%s swarm peers", ipfs)
	out, err := exec.Command("sh", "-xc", cmdLine).Output()
	if err != nil {
		return
	}
	return strings.Split(string(out), "\n")
}

func refreshPeersMap() {

	syncMap.Range(func(k, v interface{}) bool {
		syncMap.Delete(k)
		return true
	})

	for _, node := range getPeers(getIPFS()) {
		if node == "" {
			continue
		}
		syncMap.Store(node, true)
	}
}

// func swarmDisconnect2() error {
// 	sh := api.NewLocalShell()
// 	if sh == nil {
// 		return errors.New("error on connecting to local ipfs")
// 	}
// 	sh.SetTimeout(20 * time.Second)
// 	ctx := context.Background()
// 	swarm, err := sh.SwarmPeers(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	for _, node := range swarm.Peers {
// 		err = sh.Request("swarm/disconnect", node.Addr).Exec(ctx, nil)
// 	}
// 	return nil
// }

func publish(m RequestMessage) error {
	sh := api.NewLocalShell()
	if sh == nil {
		return errors.New("error on connecting to local ipfs")
	}
	sh.SetTimeout(20 * time.Second)
	msg := strings.Repeat(m.StrOption1, m.IntOption1)

	// get cid
	if msg == "" {
		return errors.New("empty string for ipfs call")
	}
	cid, err := sh.Add(strings.NewReader(msg), api.AddOpts(api.OnlyHash(true)))
	if err != nil {
		log.Println(err)
		return err
	}

	refreshPeersMap()

	// write cid to a file
	err = os.WriteFile(fmt.Sprintf("provide-%v", cid), []byte(msg), 0644)
	if err != nil {
		log.Println(err, cid)
		return err
	}

	// publish the string
	_, err = sh.Add(strings.NewReader(msg))
	if err != nil {
		log.Println(err, cid)
		return err
	}

	log.Println("publish is done:", cid)
	return nil
}

func lookup(m RequestMessage) error {
	sh := api.NewLocalShell()
	if sh == nil {
		return errors.New("error on connecting to local ipfs")
	}
	sh.SetTimeout(20 * time.Second)
	msg := strings.Repeat(m.StrOption1, m.IntOption1)

	// get cid
	if msg == "" {
		return errors.New("empty string for ipfs call")
	}
	cid, err := sh.Add(strings.NewReader(msg), api.AddOpts(api.OnlyHash(true)))
	if err != nil {
		log.Println(err)
		return err
	}

	refreshPeersMap()

	// write cid to a file
	err = os.WriteFile(fmt.Sprintf("lookup-%v", cid), []byte(msg), 0644)
	if err != nil {
		log.Println(err, cid)
		return err
	}

	_, err = sh.Cat(cid)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("lookup is done:", cid)
	return nil
}

func check(m RequestMessage) error {
	sh := api.NewLocalShell()
	if sh == nil {
		return errors.New("error on connecting to local ipfs")
	}
	sh.SetTimeout(20 * time.Second)
	msg := strings.Repeat(m.StrOption1, m.IntOption1)

	// get cid
	if msg == "" {
		return errors.New("empty string for ipfs call")
	}
	cid, err := sh.Add(strings.NewReader(msg), api.AddOpts(api.OnlyHash(true)))
	if err != nil {
		log.Println(err)
		return err
	}

	// Check
	_, err1 := os.Stat(fmt.Sprintf("ok-provide-%v", cid))
	_, err2 := os.Stat(fmt.Sprintf("ok-lookup-%v", cid))
	if err1 != nil && err2 != nil {
		log.Printf("not existing: %v", cid)
		return fmt.Errorf("not existing")
	}

	log.Println("existed.")
	os.Remove(fmt.Sprintf("ok-provide-%v", cid))
	os.Remove(fmt.Sprintf("ok-lookup-%v", cid))

	return nil
}

func getID(m RequestMessage) error {
	sh := api.NewLocalShell()
	if sh == nil {
		return errors.New("error on connecting to local ipfs")
	}
	sh.SetTimeout(20 * time.Second)
	msg := strings.Repeat(m.StrOption1, m.IntOption1)

	// get cid
	if msg == "" {
		return errors.New("empty string for ipfs call")
	}
	cid, err := sh.Add(strings.NewReader(msg), api.AddOpts(api.OnlyHash(true)))
	if err != nil {
		log.Println(err)
		return err
	}

	// Check
	_, err1 := os.Stat(fmt.Sprintf("ok-provide-%v", cid))
	_, err2 := os.Stat(fmt.Sprintf("ok-lookup-%v", cid))
	if err1 != nil && err2 != nil {
		log.Printf("not existing: %v", cid)
		return fmt.Errorf("not existing")
	}

	log.Println("existed.")
	os.Remove(fmt.Sprintf("ok-provide-%v", cid))
	os.Remove(fmt.Sprintf("ok-lookup-%v", cid))

	return nil
}
