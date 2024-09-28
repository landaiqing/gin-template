package core

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"schisandra-cloud-album/global"
	"time"
)

// InitNSQProducer 初始化生产者
func InitNSQProducer() *nsq.Producer {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(global.CONFIG.NSQ.Addr(), config)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("InitNSQ producer error: %v", err))
		return nil
	}
	return producer
}

// InitConsumer 初始化消费者
func InitConsumer(topic string, channel string) *nsq.Consumer {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 15 * time.Second
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		fmt.Printf("InitNSQ consumer error: %v\n", err)
		return nil
	}
	return consumer
}
