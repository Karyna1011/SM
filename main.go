package main

import (
	"awesomeProject/config"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
	"math/big"
	"strings"
	"time"
)

const myABI = "[{\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"set\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

func main() {
	cfg := config.NewConfig(kv.MustFromEnv())
	eth := cfg.EthClient()
	log := logan.New()

	privateKey, err := crypto.HexToECDSA(cfg.TransferConfig().Key)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	myAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	contractAddress := common.HexToAddress(cfg.TransferConfig().Address)
	chainID, err := eth.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	parsed, err := abi.JSON(strings.NewReader(myABI))
	if err != nil {
		log.Fatal("failed to parse contract ABI")
	}

	var Contract = bind.NewBoundContract(
		contractAddress,
		parsed,
		eth,
		eth,
		eth,
	)
	functionValue:=big.NewInt(200)

	d := time.NewTicker(3 * time.Second)
	for {
		select {
		case tm := <-d.C:
			nonce, err := eth.PendingNonceAt(context.Background(), myAddress )
			if err != nil {
				log.Fatal(err)
			}

			_, err = Contract.Transact(&bind.TransactOpts{
				From: contractAddress,
				Nonce:    big.NewInt(int64(nonce)),
				Signer: func(fromAddress common.Address, tx1 *types.Transaction) (*types.Transaction, error) {
					signature, err := crypto.Sign(types.NewEIP155Signer(chainID).Hash(tx1).Bytes(), privateKey)
					if err != nil {
						return nil, err
					}
					return tx1.WithSignature(types.NewEIP155Signer(chainID), signature)
				},
				Value: cfg.TransferConfig().Value,
				GasLimit: cfg.TransferConfig().GasLimit,
				GasPrice: cfg.TransferConfig().GasPrice,

			}, "set", &(functionValue))
			if err != nil {
				log.WithError(err).Error("error during calling set function")
				return
			}

			result := make([]interface{}, len(""))
			for i, e := range "" {
				result[i] = e
			}

			err = Contract.Call(&bind.CallOpts{}, &result, "get")
			if err != nil {
				log.WithError(err).Error("error during calling contract")
				return
			}

			log.Info("RESULT:", result)
			fmt.Println("The Current time is: ", tm)

		}

	}
}

