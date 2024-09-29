package mq

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"schisandra-cloud-album/core"
	"schisandra-cloud-album/global"
)

const CommentLikeTopic = "comment_like"

type MessageHandler struct{}

func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	// do whatever actual message processing is desired
	//err := processMessage(m.Body)
	fmt.Println("comment_like_mq:", string(m.Body))

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

func CommentLikeProducer(messageBody []byte) {
	err := global.NSQProducer.Publish(CommentLikeTopic, messageBody)
	if err != nil {
		global.LOG.Fatal(err)
	}
}

func CommentLikeConsumer() {
	consumer := core.InitConsumer(CommentLikeTopic)
	consumer.AddHandler(&MessageHandler{})
	err := consumer.ConnectToNSQD(global.CONFIG.NSQ.NsqAddr())
	if err != nil {
		log.Fatal(err)
	}
}
