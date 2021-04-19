package config

import (
	"load-testing/load/common"
	"runtime"
)

type LoadTestConfig struct {
	TestDuration      float32         `yaml:"duration"`
	RequestsPerSecond uint64          `yaml:"requestsPerSecond"`
	ArrivalRate       uint64          `yaml:"arrivalRate"`
	WorkersCount      int             `yaml:"workersCount"`
	LoadType          common.LoadType `yaml:"loadType"`
}

func NewLoadTestConfig() *LoadTestConfig {
	return &LoadTestConfig{
		TestDuration:      1,
		RequestsPerSecond: 1,
		WorkersCount:      runtime.GOMAXPROCS(0),
		LoadType:          common.LoadTypeDisturbed,
	}
}
