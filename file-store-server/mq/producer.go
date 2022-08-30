package mq

import (
	"fmt"
	"go_code/project14/file-store-server/config"
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 如果异常关闭，会接收通知
var notifyClose chan *amqp.Error

func init() {
	if !config.AsyncTransferEnable {
		return
	}
	if !initChannel() {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel()
			}
		}
	}()
}
func initChannel() bool {

	if channel != nil {
		return true
	}
	// init
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		fmt.Println("method initChannel failed to amqp.Dial:", err)
		return false
	}

	channel, err = conn.Channel()
	if err != nil {
		fmt.Println("method initChannel failed to conn.Channel	:", err)
		return false
	}
	return true
}

func Publish(exchange, routerkey string, msg []byte) bool {
	if !initChannel() {
		return false
	}
	// mandatory indicate  当找不到消息队列时 是否返回数据
	err := channel.Publish(exchange, routerkey, false, false, amqp.Publishing{
		Body:        msg,
		ContentType: "text/plain", // 明文 文本数据
	})
	if err != nil {
		fmt.Println("method Publish failed to channel.Publish	:", err)
		return false
	}

	return true
}
