package simplenode

import (
	"bytes"
	"context"
	"log"
	"time"

	api "github.com/ipfs/go-ipfs-api"
)

func LookupString(ctx context.Context, node string, cid string) (msg string, err error) {
	sh := api.NewShell(node)
	sh.SetTimeout(20 * time.Second)

	resp, err := sh.Cat(cid)
	if err != nil {
		log.Println(node, err)
		return

	}

	buf := &bytes.Buffer{}
	buf.ReadFrom(resp)
	msg = buf.String()
	return
}
