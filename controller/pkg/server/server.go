package server

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("agent")

// agent is the agent in the system.
type agent struct {
	ctx context.Context

	key    []byte
	server *http.Server

	handlers map[byte]func(data []byte) (byte, []byte, error)

	peerIDs []string
}

// Message type
const (
	GetID      = 1
	SetID      = 2
	Check      = 3
	Publish    = 4
	Lookup     = 5
	Clean      = 6
	Disconnect = 7
)

// NewServer creates a new server.
func NewServer(ctx context.Context, listenAddr string, keyStr string) error {
	key, err := base64.RawStdEncoding.DecodeString(keyStr)
	if err != nil {
		return err
	}
	if len(key) != 32 {
		return fmt.Errorf("Wrong key size, expect 32, got: %v", len(key))
	}
	a := &agent{
		ctx:      ctx,
		key:      key,
		handlers: make(map[byte]func(data []byte) (byte, []byte, error)),
		peerIDs:  make([]string, 0),
	}
	a.server = &http.Server{
		Addr:           listenAddr,
		Handler:        a,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// Add handlers
	a.handlers[GetID] = handleGetID
	a.handlers[SetID] = a.handleSetID
	a.handlers[Check] = handleCheck
	a.handlers[Publish] = handlePublish
	a.handlers[Lookup] = handleLookup
	a.handlers[Clean] = handleClean
	a.handlers[Disconnect] = a.handleDisonnect
	// 	Start server
	log.Infof("Start listening at %v", listenAddr)
	if err = a.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// ServeHTTP serves http.
func (a *agent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if closeErr := r.Body.Close(); closeErr != nil {
		log.Warn("HTTP can't close request body")
	}
	if err != nil {
		log.Error("Error reading request %v", err.Error())
		WriteError(w, 400, fmt.Sprintf("Invalid Request: %v.", err.Error()))
		return
	}
	if len(content) <= 1 {
		log.Error("Received content with empty request %v", content)
		WriteError(w, 400, "Content body is empty")
		return
	}
	plain, err := decrypt(content, a.key)
	if err != nil {
		log.Error("Request fails to verify")
		WriteError(w, 400, "Request fails to verify")
		return
	}
	msgType := plain[0]
	msgData := plain[1:]
	handler, ok := a.handlers[msgType]
	if !ok {
		log.Error("Unsupported message type: %v", msgType)
		WriteError(w, 400, fmt.Sprintf("Unsupported method: %v", msgType))
		return
	}
	respType, respData, err := handler(msgData)
	if err != nil {
		log.Error("Error handling request: %v", err.Error())
		WriteError(w, 400, fmt.Sprintf("Error handling request: %v", err.Error()))
		return
	}
	respEnc, err := encrypt(append([]byte{respType}, respData...), a.key)
	if err != nil {
		log.Error("Internal Error in encryption: %v", err.Error())
		WriteError(w, 500, "Internal Error")
		return
	}
	w.WriteHeader(200)
	_, err = w.Write(respEnc)
	if err != nil {
		log.Error("Error responding to client: %v", err.Error())
	}
}

// WriteError writes an error.
func WriteError(w http.ResponseWriter, header int, msg string) {
	w.WriteHeader(header)
	resp, err := json.Marshal(map[string]string{"Error": msg})
	if err == nil {
		w.Write(resp)
	}
}

func encrypt(plain []byte, key []byte) ([]byte, error) {
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	enc := aesGCM.Seal(nonce, nonce, plain, nil)
	return enc, nil
}

func decrypt(enc []byte, key []byte) ([]byte, error) {
	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	// Decrypt the data
	plain, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plain, nil
}

// shutdownRoutine shuts down the server.
func (a *agent) shutdownRoutine() {
	<-a.ctx.Done()
	a.server.Close()
}
