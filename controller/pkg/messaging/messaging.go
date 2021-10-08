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
	err = os.WriteFile(cid, []byte(msg), 0644)
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

	resp, err := sh.Cat(cid)
	if err != nil {
		log.Println(err)
		return err
	}

	_ = resp
	// buf := &bytes.Buffer{}
	// buf.ReadFrom(resp)
	// msg = buf.String()
	log.Println("lookup is done:", cid)
	return nil
}
