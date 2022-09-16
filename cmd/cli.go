package main

import (
	"flag"
	"github.com/obscuronet/faucet/faucet"
	"math/big"
)

const (
	// Flag names, defaults and usages.
	faucetPortName    = "port"
	faucetPortDefault = 80
	faucetPortUsage   = "The port on which to serve the faucet endpoint. Default: 80."

	nodeHostName    = "nodeHost"
	nodeHostDefault = "http://testnet.obscu.ro"
	nodeHostUsage   = "The host on which to connect to the Obscuro node. Default: `testnet.obscu.ro`."

	nodeWSPortName    = "nodePort"
	nodeWSPortDefault = 13000
	nodeWSPortUsage   = "The port on which to connect to the Obscuro node via RPC over HTTP. Default: 13000 ."

	faucetPKName    = "pk"
	faucetPKDefault = ""
	faucetPKUsage   = "The prefunded PK used to fund other accounts. No default, must be set."
)

func parseCLIArgs() *faucet.Config {
	faucetPort := flag.Int(faucetPortName, faucetPortDefault, faucetPortUsage)
	nodeHost := flag.String(nodeHostName, nodeHostDefault, nodeHostUsage)
	nodeWSPort := flag.Int(nodeWSPortName, nodeWSPortDefault, nodeWSPortUsage)
	faucetPK := flag.String(faucetPKName, faucetPKDefault, faucetPKUsage)
	flag.Parse()

	return &faucet.Config{
		Port:    *faucetPort,
		Host:    *nodeHost,
		WSPort:  *nodeWSPort,
		PK:      *faucetPK,
		ChainID: big.NewInt(777), // TODO make this configurable
	}
}
