
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
	"io"
	"math/big"
	"strings"
)

type Transformer interface {
	Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error)
	Reset()
}

type NopResetter struct{}

func (NopResetter) Reset() {}

type Reader struct {
	r                 io.Reader
	t                 Transformer
	err               error
	dst               []byte
	dst0, dst1        int
	src               []byte
	src0, src1        int
	transformComplete bool
}

const defaultBufSize = 4096

func NewReader(r io.Reader, t Transformer) *Reader {
	t.Reset()
	return &Reader{
		r:   r,
		t:   t,
		dst: make([]byte, defaultBufSize),
		src: make([]byte, defaultBufSize),
	}
}

const myABI = "[{\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"set\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

func main() {

	cfg := config.NewConfig(kv.MustFromEnv())
	eth := cfg.EthClient()
	log := logan.New()
	log = cfg.Log()


	privateKey, err := crypto.HexToECDSA("4f8048b22554257c143c55d3d6f56fbcdf8da0465fc0912bea0dfc44c0bf31f2")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := eth.PendingNonceAt(context.Background(), fromAddress )
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(0)
	gasLimit := uint64(300000)
	gasPrice:=big.NewInt(3000)
	toAddress := common.HexToAddress("0x118b69e0BE87a87BB30e093F496b1eE989aA15E4")
	chainID, err := eth.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	parsed, err := abi.JSON(strings.NewReader(myABI))
	if err != nil {
		log.Fatal("failed to parse contract ABI")
	}

	var Contract = bind.NewBoundContract(
		toAddress,
		parsed,
		eth,
		eth,
		eth,
	)
	result := make([]interface{}, len(""))
	for i, e := range "" {
		result[i] = e
	}

	functionValue:=big.NewInt(42)

	_, err = Contract.Transact(&bind.TransactOpts{
		From: toAddress,
		Nonce:    big.NewInt(int64(nonce)),
		Signer: func(fromAddress common.Address, tx1 *types.Transaction) (*types.Transaction, error) {
			signature, err := crypto.Sign(types.NewEIP155Signer(chainID).Hash(tx1).Bytes(), privateKey)
			if err != nil {
				return nil, err
			}
			return tx1.WithSignature(types.NewEIP155Signer(chainID), signature)
		},
		Value: value,
		GasLimit: gasLimit,
		GasPrice: gasPrice,


	}, "set", &(functionValue))
	if err != nil {
		log.WithError(err).Error("error during calling set function")
		return
	}

	err = Contract.Call(&bind.CallOpts{}, &result, "get")
	if err != nil {
		log.WithError(err).Error("error during calling contract")
		return
	}

	fmt.Println("RESULT:", result)
}
