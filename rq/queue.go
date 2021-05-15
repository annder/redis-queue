package rq

import (
	"context"
	"goredisqueue/msg"
	"log"

	"github.com/go-redis/redis/v8"
)

type Queue struct {
	con *redis.Client
}

var ctx = context.TODO()

// 删除队列中的消息
func (q *Queue) lrem(queue, msg string) error {

	defer log.Fatalln(q.con.Close())

	if _, err := q.con.LRem(ctx, queue, 1, msg).Result(); err != nil {
		return err
	}
	return nil
}

func (q *Queue) rpoplpush(imsg msg.IMessage, sourceQueue, destQueue string) (interface{}, msg.IMessage, error) {
	//  把数据放在最前面
	r, err := q.con.RPopLPush(ctx, sourceQueue, destQueue).Result() // 这里要改一下结构

	if err != nil {
		return nil, nil, err
	}
	if r == "" {
		return nil, nil, nil
	}
	defer q.con.Close()

}

// 接受数据
func (q *Queue) Receive(queue, msg string) string {
	result, _ := q.con.RPopLPush(ctx, queue, msg).Result()
	ok := q.lrem(queue, msg)
	log.Println(ok, result)
	if ok == nil {
		return result
	}
	return ""
}

// 传递数据
func (q *Queue) Delivery(queue, msg string) bool {
	str, _ := q.con.LPush(ctx, queue, msg).Result()
	if str > 0 {
		return true
	}
	return false
}

// 新建一个连接
func NewCon(instance *redis.Client) *Queue {
	return &Queue{
		con: instance,
	}
}
