package main

import (
	"context"
	"fmt"
	"load-testing/config"
	"load-testing/core/dispatcher"
	"load-testing/load"
	"os"
)

var (
	A int = 0
)

func prinSmth() error {
	A++

	return nil
}

func main() {
	cfg, err := config.ParseCLIParams(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	disp := dispatcher.NewRoundRobinDispatcher(*cfg)

	ls := load.NewLoadService(disp, context.Background())
	ls.SetLoadTime(cfg.TestDuration)
	ls.AddJob(prinSmth)

	ls.Start()
	
	fmt.Println(A)
}
