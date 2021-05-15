package rq

import (
	"context"
	"errors"
	"fmt"
	"goredisqueue/msg"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Queue struct {
	con *redis.Client
}

var ctx = context.TODO()

// 删除队列中的消息
func (q *Queue) lrem(queue, msg string) error {

	if _, err := q.con.LRem(ctx, queue, 1, msg).Result(); err != nil {
		return err
	}
	return nil
}

func (q *Queue) rpoplpush(imsg msg.IMessage, sourceQueue, destQueue string) (interface{}, msg.IMessage, error) {
	//  TODO，这里的问题
	r, err := q.con.Do(ctx, "RPOPLPUSH", sourceQueue, destQueue).Result()
	log.Println("rpoplpush -> ", r)
	if err != nil {
		return nil, nil, err
	}
	if r == nil {
		return nil, nil, nil
	}
	rUint8, ok := r.([]uint8)
	//  如果不是uint8的话
	if !ok {
		return nil, nil, errors.New("error")
	}
	if msg_, err := imsg.Unmarshal(rUint8); err != nil {
		return nil, nil, err
	} else if _, ok := msg_.(msg.IMessage); ok {
		//  是否实现了接口
		return r, msg_, nil
	} else {
		// 无法实现接口
		return nil, nil, errors.New("cannot assert msg as interface IMessage")
	}
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
func (q *Queue) Delivery(msg msg.IMessage) error {
	// 之前的消息

	perpareMsg := fmt.Sprintf("%s.prepare", msg.GetChannel())
	if toMsgJSON, err := msg.Marshal(); err != nil {
		return err
	} else {
		_, err := q.con.LPush(ctx, perpareMsg, toMsgJSON).Result()
		return err
	}
}

// 初始化接受值
func (q *Queue) InitReceiver(msg msg.IMessage) {
	// 之前的队列名
	prepareQueue := fmt.Sprintf("%s.prepare", msg.GetChannel())
	doingQueue := fmt.Sprintf("%s.doing", msg.GetChannel())
	// 这块有问题
	go func() {
		for {
			reply, msg, err := q.rpoplpush(msg, prepareQueue, doingQueue)
			toStringReplay, ok := reply.(string)
			log.Println("push ->", toStringReplay)
			if !ok {
				errors.New("assert replay is failure, becuase it not string type.")
			}
			// 这里有问题
			if err != nil {
				log.Println("queque -> 93 ")
				log.Println(err)

			}
			// 如果是空的数据
			if msg == nil {
				continue
			}
			if err := msg.Resolve(); err == nil {
				_ = q.lrem(doingQueue, toStringReplay)
				log.Println("消费了")
			} else {
				log.Fatalln(err)
			}
		}
	}()
	fmt.Printf("receiver have been initialized\n")
}

func (q *Queue) SetSomething(str string) {
	ok, _ := q.con.Set(ctx, "key1", str, time.Hour).Result()
	log.Println(ok)
}

// 新建一个连接
func NewCon(instance *redis.Client) *Queue {
	return &Queue{
		con: instance,
	}
}
