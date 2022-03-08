package ws

import (
	"encoding/json"
	"im/ay"
	"im/models"
	"im/models/api"
	"log"
	"time"
)

// 私聊

type Chat struct {
}

func (con Chat) Operate(msg Message, data []byte) {

	cc, err := json.Marshal(&api.Msg{
		SendId:    msg.SendId,
		ReceiveId: msg.ReceiveId,
		Content:   msg.Content,
		Type:      msg.Type,
		//CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		CreatedAt: models.JsonTime(time.Now()),
	})

	_, err = ay.Redis.Do("lpush", "message", cc)
	if err != nil {
		log.Println(err)
	}
	SendMsg(msg.ReceiveId, data)
}
