package main

import (
	"goredisqueue/client"
	"goredisqueue/rq"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	redisClient := client.Reids()
	queue := rq.NewCon(redisClient)
	s := gin.Default()
	// 发送队列数据
	s.GET("/delivery", func(context *gin.Context) {
		if ok := queue.Delivery("akv", "msg"); ok {
			context.JSON(200, gin.H{
				"code": 200,
				"msg":  "ok",
				"data": nil,
			})
			return
		} else {
			context.JSON(200, gin.H{
				"code": 500,
				"msg":  "err",
				"data": nil,
			})
			context.Abort()
			return
		}
	})
	// 接受队列数据
	s.GET("/receive", func(context *gin.Context) {
		value := queue.Receive("akv", "msg")
		if value != "" {
			context.JSON(200, gin.H{
				"code": 200,
				"msg":  "ok",
				"data": value,
			})
			return
		}
	})
	log.Fatal(s.Run(":8080"))
}
