package faucet

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"sync"
	"time"

	"math/big"
)

const _timeout = 30 * time.Second

type Faucet struct {
	client    *obsclient.AuthObsClient
	fundMutex sync.Mutex
	nonce     uint64
	wallet    wallet.Wallet
}

func (f *Faucet) Fund(address *common.Address) error {
	// the faucet should be the only user of the faucet pk

	// todo remove hardcoded gas values
	gas := uint64(21000)

	f.fundMutex.Lock()
	tx := &types.LegacyTx{
		Nonce:    f.nonce,
		GasPrice: big.NewInt(225),
		Gas:      gas,
		To:       address,
		Value:    new(big.Int).Mul(big.NewInt(10), big.NewInt(params.Ether)),
	}

	signedTx, err := f.wallet.SignTransaction(tx)
	if err != nil {
		f.fundMutex.Unlock()
		return err
	}

	if err := f.client.SendTransaction(context.Background(), signedTx); err != nil {
		f.fundMutex.Unlock()
		return err
	}
	f.nonce++
	f.fundMutex.Unlock()

	txMarshal, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	fmt.Printf("Funded address: %s - tx: %+v\n", address.Hex(), string(txMarshal))
	// todo handle tx receipt

	if err := f.validateTx(signedTx); err != nil {
		return fmt.Errorf("unable to validate tx %s: %w", signedTx.Hash(), err)
	}

	return nil
}

func (f *Faucet) validateTx(tx *types.Transaction) error {
	for now := time.Now(); time.Since(now) < _timeout; time.Sleep(time.Second) {
		receipt, err := f.client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			if errors.Is(err, rpcclientlib.ErrNilResponse) {
				// tx receipt is not available yet
				continue
			}
			return fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
		}

		txReceiptBytes, err := receipt.MarshalJSON()
		if err != nil {
			return fmt.Errorf("could not marshal transaction receipt to JSON in eth_getTransactionReceipt request. Cause: %w", err)
		}
		fmt.Println(string(txReceiptBytes))

		if receipt.Status != 1 {
			return fmt.Errorf("tx status is not 0x1")
		}
		return nil
	}
	return fmt.Errorf("unable to fetch tx receipt after %s", _timeout)
}

func NewFaucet(rpcUrl string, chainID *big.Int, pk *ecdsa.PrivateKey) (*Faucet, error) {
	w := wallet.NewInMemoryWalletFromPK(chainID, pk)
	obsClient, err := obsclient.DialWithAuth(rpcUrl, w)
	if err != nil {
		return nil, fmt.Errorf("unable to connect with the node: %w", err)
	}

	nonce, err := obsClient.NonceAt(context.Background(), w.Address(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch %s nonce: %w", w.Address(), err)
	}
	return &Faucet{
		client: obsClient,
		wallet: w,
		nonce:  nonce,
	}, nil
}
