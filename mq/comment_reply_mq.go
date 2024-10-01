package mq

import (
	"github.com/nsqio/go-nsq"
	"log"
	"schisandra-cloud-album/core"
	"schisandra-cloud-album/global"
)

const CommentReplyTopic = "comment_reply"

type CommentReplyMessageHandler struct{}

// CommentReplyProducer 评论回复消息生产
func CommentReplyProducer(messageBody []byte) {
	err := global.NSQProducer.Publish(CommentReplyTopic, messageBody)
	if err != nil {
		global.LOG.Fatal(err)
	}
}

// CommentReplyConsumer 评论回复消息消费
func CommentReplyConsumer() {
	consumer := core.InitConsumer(CommentReplyTopic)
	consumer.AddHandler(&CommentReplyMessageHandler{})
	err := consumer.ConnectToNSQD(global.CONFIG.NSQ.NsqAddr())
	if err != nil {
		log.Fatal(err)
	}
}

func (h *CommentReplyMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	return nil
}
