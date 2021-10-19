package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
)

type config struct {
	transferConfig TransferConfig
	contractConfig ContractConfig
	getter         kv.Getter
	onceTransfer   comfig.Once
	onceContract   comfig.Once

	Ether
	comfig.Logger
}

type Config interface {
	TransferConfig() TransferConfig
	ContractConfig() ContractConfig
	Log() *logan.Entry
	Ether
}

func NewConfig(getter kv.Getter) Config {
	return &config{
		getter: getter,
		Ether:  NewEther(getter),
		Logger: comfig.NewLogger(getter, comfig.LoggerOpts{Release: "0.1.0"}),
	}
}

