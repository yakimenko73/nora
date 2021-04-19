package config

import (
	"github.com/jessevdk/go-flags"
	"load-testing/load/common"
)

type CLIConfig struct {
	TestDuration      float32 `short:"d" long:"duration" description:"total test duration (in minutes)" required:"true"`
	RequestsPerSecond uint64  `short:"r" long:"rps" description:"requests/second" required:"true"`
	WorkersCount      int     `short:"w" long:"workers-count" description:"number of workers that will execute jobs" required:"true"`
	ArrivalRate       uint64  `short:"a" long:"arrival-rate" description:"arrival rate" required:"false" default:"0"`
	LoadType          string  `short:"t" long:"load-type" description:"load type" required:"false" default:"disturbed" choice:"ramp" choice:"disturbed"`
}

func ParseCLIParams(args []string) (*LoadTestConfig, error) {
	var cliConfig CLIConfig

	_, err := flags.ParseArgs(&cliConfig, args)
	if err != nil {
		return nil, err
	}

	return parseCLIConfigToLoadTestConfig(&cliConfig), nil
}

func parseCLIConfigToLoadTestConfig(cliConfig *CLIConfig) *LoadTestConfig {
	return &LoadTestConfig{
		TestDuration:      cliConfig.TestDuration,
		RequestsPerSecond: cliConfig.RequestsPerSecond,
		ArrivalRate:       cliConfig.ArrivalRate,
		WorkersCount:      cliConfig.WorkersCount,
		LoadType:          common.LoadType(cliConfig.LoadType),
	}
}
