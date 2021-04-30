package main

import (
	"context"
	"fmt"
	"load-testing/config"
	"load-testing/core/dispatcher"
	"load-testing/core/executor"
	"load-testing/load"
	"os"
	"time"
)

var (
	TOTAL = 1

	//conn *amqp.Connection
	//channel *amqp.Channel
)


func prinSmth() error {
	TOTAL++

	//return channel.Publish(
	//	"amq.direct",
	//	"r.test",
	//	false,
	//	false,
	//	amqp.Publishing{
	//		Body: []byte("privet"),
	//	},
	//)
	return nil
}

func init() {
	//var err error
	//conn, err = amqp.Dial("amqp://admin:admin@localhost:5772/")
	//if err != nil {
	//	panic(err)
	//}
	//
	//channel, err = conn.Channel()
	//if err != nil {
	//	panic(err)
	//}
}

func main() {
	cfg, err := config.ParseCLIParams(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	exec := executor.New(cfg)
	disp := dispatcher.NewRoundRobinDispatcher(*cfg, exec)

	ls := load.NewLoadService(disp, context.Background())
	ls.SetLoadTime(cfg.TestDuration)
	_ = ls.AddJob(prinSmth)

	fmt.Println(time.Now())
	ls.Start()
	fmt.Println(time.Now())

	fmt.Println(float32(TOTAL) / 1000000)
}
