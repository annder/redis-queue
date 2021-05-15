package main

import (
	"goredisqueue/client"
	msg2 "goredisqueue/msg"
	"goredisqueue/rq"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	redisClient := client.Reids()
	queue := rq.NewCon(redisClient)
	msg := &msg2.Message{Name: "domeQueue"}
	queue.InitReceiver(msg)
	go func() {
		for i := 0; i < 10; i++ {
			msg := &msg2.Message{Name: "demoQueue", Content: map[string]string{
				"order_no": strconv.FormatInt(time.Now().Unix(), 10),
			}}
			_ = queue.Delivery(msg)
		}
	}()
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT)
	//
	for {

		switch <-quit {
		case syscall.SIGINT:
			return
		}
	}
}
