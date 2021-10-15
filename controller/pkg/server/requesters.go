package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Request sends a request to given addr with given key, msg type and data.
func Request(addr string, key []byte, msgType byte, data []byte) (byte, []byte, error) {
	if !strings.HasPrefix(addr, "http://") {
		addr = "http://" + addr
	}
	enc, err := encrypt(append([]byte{msgType}, data...), key)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", addr, bytes.NewReader(enc))
	if err != nil {
		return 0, nil, err
	}
	client := &http.Client{Timeout: 90 * time.Second}
	r, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	content, err := ioutil.ReadAll(r.Body)
	if closeErr := r.Body.Close(); closeErr != nil {
		return 0, nil, closeErr
	}
	if len(content) <= 1 {
		return 0, nil, fmt.Errorf("Received content with empty request %v", content)
	}
	plain, err := decrypt(content, key)
	if err != nil {
		return 0, nil, err
	}
	return plain[0], plain[1:], nil
}

// RequestGetID gets peer id.
func RequestGetID(addr string, key []byte) (string, error) {
	msgType, data, err := Request(addr, key, GetID, []byte{0})
	if err != nil {
		return "", err
	}
	if msgType != GetID {
		return "", fmt.Errorf("wrong msg type received, expect %v, got %v", GetID, msgType)
	}
	return string(data), nil
}

// RequestSetID sets peer ids.
func RequestSetID(addr string, key []byte, ids []string) ([]string, error) {
	msgType, data, err := Request(addr, key, SetID, []byte(strings.Join(ids, ";")))
	if err != nil {
		return nil, err
	}
	if msgType != SetID {
		return nil, fmt.Errorf("wrong msg type received, expect %v, got %v", SetID, msgType)
	}
	return strings.Split(string(data), ";"), nil
}

// RequestCheck requests check.
func RequestCheck(addr string, key []byte, cid string) (bool, error) {
	msgType, data, err := Request(addr, key, Check, []byte(cid))
	if err != nil {
		return false, err
	}
	if msgType != Check {
		return false, fmt.Errorf("wrong msg type received, expect %v, got %v", Check, msgType)
	}
	if data[0] == 1 {
		return false, nil
	}
	return true, nil
}

// RequestPublish requests publish a random content.
func RequestPublish(addr string, key []byte, content []byte) (string, error) {
	msgType, data, err := Request(addr, key, Publish, content)
	if err != nil {
		return "", err
	}
	if msgType != Publish {
		return "", fmt.Errorf("wrong msg type received, expect %v, got %v", Publish, msgType)
	}
	return string(data), nil
}

// RequestLookup requests a lookup.
func RequestLookup(addr string, key []byte, cid string) error {
	msgType, data, err := Request(addr, key, Lookup, []byte(cid))
	if err != nil {
		return err
	}
	if msgType != Lookup {
		return fmt.Errorf("wrong msg type received, expect %v, got %v", Publish, msgType)
	}
	if cid != string(data) {
		return fmt.Errorf("look up the wrong cid, expect %v, got %v", cid, string(data))
	}
	return nil
}

// RequestClean requests a clean.
func RequestClean(addr string, key []byte, cid string) (string, error) {
	msgType, data, err := Request(addr, key, Clean, []byte(cid))
	if err != nil {
		return "", err
	}
	if msgType != Clean {
		return "", fmt.Errorf("wrong msg type received, expect %v, got %v", Clean, msgType)
	}
	return string(data), nil
}

// RequestDisconnect requests a disconnection.
func RequestDisconnect(addr string, key []byte) (string, error) {
	msgType, data, err := Request(addr, key, Disconnect, []byte{1})
	if err != nil {
		return "", err
	}
	if msgType != Disconnect {
		return "", fmt.Errorf("wrong msg type received, expect %v, got %v", Disconnect, msgType)
	}
	return string(data), nil
}
