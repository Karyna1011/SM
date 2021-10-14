package config

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Ether interface {
	EthClient() *ethclient.Client
}

type ether struct {
	getter kv.Getter
	once   comfig.Once
}

func NewEther(getter kv.Getter) Ether {
	return &ether{getter: getter}
}

func (h *ether) EthClient() *ethclient.Client {
	return h.once.Do(func() interface{} {
		var config struct {
			Endpoint string `fig:"endpoint,required"`
		}

		err := figure.
			Out(&config).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(h.getter, "rpc")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out rpc"))
		}

		eth, err := ethclient.Dial(config.Endpoint)
		if err != nil {
			panic(fmt.Sprintf("failed to dial %s", config.Endpoint))
		}

		return eth
	}).(*ethclient.Client)
}
func (h *ether) signer() *ethclient.Client {
	return h.once.Do(func() interface{} {
		var config struct {
			signer string `fig:"signer,required"`
		}

		err := figure.
			Out(&config).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(h.getter, "withdraw")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out contractData"))
		}

		fmt.Printf("Ethereumt withdraw: %s", config.signer)

		eth, err := ethclient.Dial(config.signer)
		if err != nil{
			panic(fmt.Sprintf("failed to dial %s", config.signer))
		}

		return eth

		return 0
	}).(*ethclient.Client)
}

