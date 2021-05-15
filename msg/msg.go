package msg

import (
	"context"
	"encoding/json"
	"fmt"
)

var ctx = context.TODO()

type IMessage interface {
	Resolve() error
	GetChannel() string // 获取通道
	Marshal() ([]byte, error)
	Unmarshal([]uint8) (IMessage, error)
}

type Message struct {
	Name    string            // 消息名
	Content map[string]string `json:"content"`
}

// 获取通道
func (m *Message) GetChannel() string {
	return m.Name
}

// 序列化
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}


func (m *Message) Resolve()  error {
	// 做一个逻辑处理
	fmt.Printf("consumed %+v\n", m.Content)
	return nil
}

//反序列化
func (m *Message) Unmarshal(replay []byte) (IMessage, error) {
	var msg Message
	err := json.Unmarshal(replay, &msg)
	return &msg, err
}

