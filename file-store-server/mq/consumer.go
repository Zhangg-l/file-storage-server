package mq

import "fmt"

var done chan struct{}

func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	// 1.通过channel获取消息
	msgCh, err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println("in StartConsume method failed to  channel.Consume: ", err)
		return
	}
	// listen channel
	done = make(chan struct{}, 1)
	go func() {
		for msg := range msgCh {
			fmt.Println("msg is processing")
			if !callback(msg.Body) {
				// todo: 错误处理
			}
		}
	}()
	fmt.Println("done...")
	<-done
	fmt.Println("done...")
	channel.Close()
}
