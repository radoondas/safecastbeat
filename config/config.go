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
	//Measurements: []Measurement{
	//	{
	//		Type:              "added",
	//		SingleRun:         false,
	//		Since:             "",
	//		Until:             "",
	//		CapturedAfter:     "",
	//		CapturedBefore:    "",
	//		Latitude:          0,
	//		Longitude:         0,
	//		Distance:          0,
	//		SingleRunFinished: false,
	//	},
	//},
	// need to declare default measurment??
}

type SSL struct {
	CAfile string `config:"cafile"`
}

//type Measurement struct {
//	Type           string  `config:"type"`
//	SingleRun      bool    `config:"single_run"`
//	Since          string  `config:"since"`
//	Until          string  `config:"until"`
//	CapturedAfter  string  `config:"captured_after"`
//	CapturedBefore string  `config:"captured_before"`
//	Latitude       float64 `config:"latitude"`
//	Longitude      float64 `config:"longitude"`
//	Distance       float64 `config:"distance"`
//	SingleRunFinished bool //default is false which works for us
//}
