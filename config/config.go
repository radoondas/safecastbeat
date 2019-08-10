// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"time"
)

type Config struct {
	Period time.Duration `config:"period"`
	SSL    SSL           `config:"ssl"`
	//Measurements []Measurement `config:"measurements"`
}

var DefaultConfig = Config{
	Period: 10 * time.Second,
	SSL: SSL{
		CAfile: "",
	},
}

type SSL struct {
	CAfile string `config:"cafile"`
}
