package ws

import (
	"encoding/json"
	"im/ay"
	"log"
	"time"
)

// 私聊

type Chat struct {
}

func (con Chat) Operate(msg Message, data []byte) {

	isRead := 1

	cc, err := json.Marshal(&Message{
		SendId:    msg.SendId,
		ReceiveId: msg.ReceiveId,
		Content:   msg.Content,
		Type:      msg.Type,
		IsRead:    isRead,
		SendAt:    msg.CreatedAt,
		CreatedAt: time.Now().Unix(),
	})

	_, err = ay.Redis.Do("lpush", "message", cc)
	if err != nil {
		log.Println("operate:" + err.Error())
	}
	SendMsg(msg.ReceiveId, data)
}
