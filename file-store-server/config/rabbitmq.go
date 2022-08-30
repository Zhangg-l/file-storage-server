package config

const (
	AsyncTransferEnable    = true
	RabbitURL              = "amqp://guest:guest@127.0.0.1:5672/"
	TransExchangeName      = "uploadserver.trans"
	TransOSSQueueName      = "uploadserver.trans.oss"
	TransOSSErrQueueName   = "uploadserver.trans.oss.err"
	TransOSSRountingKey    = "oss"
)
