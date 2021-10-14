package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math/big"
)

type TransferConfig struct {
	Key           string   `fig:"key"`
	Address       string   `fig:"address"`
	Value         *big.Int `fig:"value"`
	GasLimit      uint64   `fig:"gas_limit"`
	GasPrice      *big.Int `fig:"gas_price"`
}

func (c *config) TransferConfig() TransferConfig {
	c.once.Do(func() interface{} {
		var result TransferConfig

		err := figure.Out(&result).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "transfer")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out transfer"))
		}
		c.transferConfig = result
		return nil
	})
	return c.transferConfig
}

