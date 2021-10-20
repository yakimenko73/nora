package task_scheduler

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSandbox(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tsch, err := New(
		WithConsoleUserInterface(),
		WithDuration(60*time.Second),
		WithTask("rnd task", rndTask),
	)
	if err != nil {
		t.Fatal(err)
	}

	for metric := range tsch.Run(ctx) {
		fmt.Printf("metric - %v\n", metric)
	}
}

func rndTask(ctx context.Context) error {
	time.Sleep(time.Duration(int(time.Second)*(rand.Int()%3) + 1))
	return nil
}

//func TestSandbox(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//
//	main, err := screen.NewMainScreen()
//
//	ter, err := createTerminal()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	c, err := cui.NewCui(ter, main)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	done := make(chan bool)
//	done2 := make(chan bool)
//	go func() {
//		err := c.Run(ctx, done)
//		fmt.Println(err)
//		done2 <- true
//	}()
//
//	<- done
//	<- done2
//}
