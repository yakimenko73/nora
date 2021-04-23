package main

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"load-testing/config"
	"load-testing/core/dispatcher"
	"load-testing/load"
	"os"
	"time"
)

var (
	A int = 0

	conn *amqp.Connection
	channel *amqp.Channel
)

func prinSmth() error {
	A++

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
	var err error
	conn, err = amqp.Dial("amqp://admin:admin@localhost:5772/")
	if err != nil {
		panic(err)
	}

	channel, err = conn.Channel()
	if err != nil {
		panic(err)
	}
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

	fmt.Println(time.Now())
	ls.Start()
	fmt.Println(time.Now())

	fmt.Println(float32(A) / 1000000)
}
