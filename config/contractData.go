package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type ContractConfig struct {
	AmountOutMin  int64       `fig:"amountOutMin"`
	AddressArray  []string    `fig:"addressArray"`
	Address       string      `fig:"address"`
	Deadline      int64       `fig:"deadline"`
}

func (c *config) ContractConfig() ContractConfig {
	c.once.Do(func() interface{} {
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

