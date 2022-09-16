package faucet

import "math/big"

type Config struct {
	Port     int
	Host     string
	HTTPPort int
	PK       string
	ChainID  *big.Int
}
