package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/faucet/faucet"
)

func main() {
	cfg := parseCLIArgs()

	if cfg.PK == "" {
		panic("no key loaded")
	}
	nodeAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.WSPort)
	key, err := crypto.HexToECDSA(cfg.PK[2:])
	if err != nil {
		panic(err)
	}

	f, err := faucet.NewFaucet(nodeAddr, cfg.ChainID, key)
	if err != nil {
		panic(err)
	}
	server := faucet.NewWebServer(f)
	server.Start()
}
