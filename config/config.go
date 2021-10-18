package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math/big"
	"time"
)

type TransferConfig struct {
	Key       string        `fig:"key"`
	Address   string        `fig:"address"`
	GasLimit  uint64        `fig:"gas_limit"`
	GasPrice  *big.Int      `fig:"gas_price"`
	Time      time.Duration `fig:"time"`
	Timestamp time.Duration `fig:"timestamp"`
}

func (c *config) TransferConfig() TransferConfig {
	c.onceTransfer.Do(func() interface{} {
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


