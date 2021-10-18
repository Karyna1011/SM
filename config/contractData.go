package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math/big"
)

type ContractConfig struct {
	Percent       *big.Int     `fig:"percent"`
	Amount        *big.Int    `fig:"amount"`
	AddressArray  []string    `fig:"addresses"`
}

func (c *config) ContractConfig() ContractConfig {
	c.onceContract.Do(func() interface{} {
		var result ContractConfig

		err := figure.Out(&result).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "contractData")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out contractData"))
		}
		c.contractConfig = result
		return nil
	})
	return c.contractConfig
}

