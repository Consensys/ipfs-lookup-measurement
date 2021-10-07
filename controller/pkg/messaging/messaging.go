package messaging

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	api "github.com/ipfs/go-ipfs-api"
)

type RequestMessage struct {
	Cmd        string
	IntOption1 int
	StrOption1 string
}

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
			err := json.NewDecoder(r.Body).Decode(&m)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "failed")
				log.Println(err)
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

func taskHandler(m RequestMessage) error {
	switch m.Cmd {
	case "publish":
		return publish(m)
	case "lookup":
		return lookup(m)
	}
	return errors.New("command is invalid")
}

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

	// write cid to a file
	err = cidToFile(cid, msg)
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

	resp, err := sh.Cat(cid)
	if err != nil {
		log.Println(err)
		return err
	}

	_ = resp
	// buf := &bytes.Buffer{}
	// buf.ReadFrom(resp)
	// msg = buf.String()
	return nil
}

func cidToFile(cid string, msg string) error {
	ipfsTestFolder := os.Getenv("PERFORMANCE_TEST_DIR")
	if ipfsTestFolder == "" {
		ipfsTestFolder = "ipfs-test"
	}
	err := os.RemoveAll(ipfsTestFolder)
	if err != nil {
		log.Println(err, ipfsTestFolder)
		return err
	}
	err = os.MkdirAll(ipfsTestFolder, 0755)
	if err != nil {
		log.Println(err, ipfsTestFolder)
		return err
	}
	return os.WriteFile(ipfsTestFolder+"/"+cid, []byte(msg), 0644)
}
