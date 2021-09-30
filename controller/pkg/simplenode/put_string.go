package simplenode

import (
	"context"
	"log"
	"strings"

	api "github.com/ipfs/go-ipfs-api"
)

func PutString(ctx context.Context, node string, msg string) (cid string, err error) {
	sh := api.NewShell(node)

	cid, err = sh.Add(strings.NewReader(msg))
	if err != nil {
		log.Println(node, err)
		return
	}
	return
}

func GetCid(ctx context.Context, node string, msg string) (cid string) {
	sh := api.NewShell(node)

	cid, err := sh.Add(strings.NewReader(msg), api.AddOpts(api.OnlyHash(true)))
	if err != nil {
		log.Println(node, err)
		return
	}

	return
}
