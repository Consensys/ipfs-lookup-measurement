package simplenode

import (
	"bytes"
	"context"
	"log"

	api "github.com/ipfs/go-ipfs-api"
)

func LookupString(ctx context.Context, node string, cid string) (msg string, err error) {
	sh := api.NewShell(node)

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
