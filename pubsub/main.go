package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/libp2p/go-floodsub"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs-addr"

	"github.com/multiformats/go-multihash"
)

func main() {
	ctx := context.Background()
	host, err := libp2p.New(ctx, libp2p.Defaults)
	if err != nil {
		panic(err)
	}
}

fsub, err := floodsub.NewFloodSub(ctx, host)
if err != nil {
panic(err)
}

