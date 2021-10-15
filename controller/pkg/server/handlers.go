package server

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	api "github.com/ipfs/go-ipfs-api"
)

// getIPFSCLI gets the ipfs cli
func getIPFSCLI() string {
	ipfs := os.Getenv("IPFS")
	if ipfs == "" {
		ipfs = "/home/ubuntu/go-ipfs/cmd/ipfs/ipfs"
	}
	return ipfs
}

// handleGetID handles get id request.
func handleGetID(data []byte) (byte, []byte, error) {
	// Use IPFS shell
	sh := api.NewLocalShell()
	if sh == nil {
		return 0, nil, fmt.Errorf("error getting local ipfs shell")
	}
	sh.SetTimeout(20 * time.Second)

	id, err := sh.ID()
	if err != nil {
		return 0, nil, err
	}
	return GetID, []byte(id.ID), nil
}

// handleSetID handles set id request.
func (a *agent) handleSetID(data []byte) (byte, []byte, error) {
	idStr := string(data)
	idStrs := strings.Split(idStr, ";")
	for _, id := range idStrs {
		a.peerIDs = append(a.peerIDs, id)
	}
	return SetID, []byte(strings.Join(a.peerIDs, ";")), nil
}

// handleCheck handles check request.
func handleCheck(data []byte) (byte, []byte, error) {
	// Use IPFS shell
	sh := api.NewLocalShell()
	if sh == nil {
		return 0, nil, fmt.Errorf("error getting local ipfs shell")
	}
	sh.SetTimeout(20 * time.Second)

	// Get cid
	cid := string(data)

	// Check
	_, err1 := os.Stat(fmt.Sprintf("ok-provide-%v", cid))
	_, err2 := os.Stat(fmt.Sprintf("ok-lookup-%v", cid))
	if err1 != nil && err2 != nil {
		log.Infof("not existing: %v\n", cid)
		return Check, []byte{1}, nil
	}
	log.Infof("existing: %v\n", cid)
	os.Remove(fmt.Sprintf("ok-provide-%v", cid))
	os.Remove(fmt.Sprintf("ok-lookup-%v", cid))

	return Check, []byte{0}, nil
}

// handlePublish handles publish request.
func handlePublish(data []byte) (byte, []byte, error) {
	// Use IPFS shell
	sh := api.NewLocalShell()
	if sh == nil {
		return 0, nil, fmt.Errorf("error getting local ipfs shell")
	}
	sh.SetTimeout(20 * time.Second)

	// Get cid
	cid, err := sh.Add(bytes.NewReader(data), api.AddOpts(api.OnlyHash(true)))
	if err != nil {
		return 0, nil, err
	}

	// write cid to a file
	err = os.WriteFile(fmt.Sprintf("provide-%v", cid), []byte{1}, 0644)
	if err != nil {
		return 0, nil, fmt.Errorf("error writing cid %v to file: %v\n", cid, err)
	}

	// Publish the content
	_, err = sh.Add(bytes.NewReader(data))
	if err != nil {
		return 0, nil, err
	}

	return Publish, []byte(cid), nil
}

// handleLookup handles lookup request.
func handleLookup(data []byte) (byte, []byte, error) {
	// Use IPFS shell
	sh := api.NewLocalShell()
	if sh == nil {
		return 0, nil, fmt.Errorf("error getting local ipfs shell")
	}
	sh.SetTimeout(20 * time.Second)

	// Get cid
	cid := string(data)

	// write cid to a file
	err := os.WriteFile(fmt.Sprintf("lookup-%v", cid), []byte{1}, 0644)
	if err != nil {
		return 0, nil, fmt.Errorf("error writing cid %v to file: %v\n", cid, err)
	}

	// Retrieve the content
	_, err = sh.Cat(cid)
	if err != nil {
		return 0, nil, err
	}

	// Add log
	f, err := os.OpenFile("/home/ubuntu/all.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, nil, err
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("Retrieval for %v is finished.\n", cid)); err != nil {
		return 0, nil, err
	}

	return Lookup, []byte(cid), nil
}

// handleClean handles clean request.
func handleClean(data []byte) (byte, []byte, error) {
	// Use IPFS CLI
	ipfs := getIPFSCLI()

	// Get cid
	cid := string(data)

	cli := fmt.Sprintf("%s pin rm %s; %s repo gc", ipfs, cid, ipfs)
	out, err := exec.Command("sh", "-xc", cli).CombinedOutput()
	if err != nil {
		return 0, nil, err
	}

	return Clean, out, nil
}

// handleDisconnect handles disconnect request.
func (a *agent) handleDisonnect(data []byte) (byte, []byte, error) {
	// Use IPFS CLI
	ipfs := getIPFSCLI()

	output := make([]string, 0)
	for _, peer := range a.peerIDs {
		cli := fmt.Sprintf("%v swarm peers | /bin/grep %s | %v swarm disconnect", ipfs, peer, ipfs)
		out, err := exec.Command("sh", "-xc", cli).CombinedOutput()
		if err != nil {
			output = append(output, fmt.Sprintf("%v%v\n", string(out), err.Error()))
		} else {
			output = append(output, string(out))
		}
	}

	return Disconnect, []byte(strings.Join(output, ";")), nil
}
