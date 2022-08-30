package faucet

import "math/big"

type Config struct {
	Port    int
	Host    string
	WSPort  int
	PK      string
	ChainID *big.Int
}
