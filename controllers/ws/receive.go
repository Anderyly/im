package ws

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"im/ay"
	"im/models"
	"im/models/api"
	"log"
	"time"
)

func Receive(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			ay.Redis.Do("ZREM", "online_users", node.Id)
			node.Conn.Close()
			log.Println("Receive:" + err.Error())
			return
		}

		Dispatch(node.Id, data)
		log.Printf("recv<=%s", data)
	}
}

func MessageToMysql() {
	for {
		res, err := redis.String(ay.Redis.Do("rpop", "message"))
		if err != nil {
			log.Println("当前未查询到redis中消息")
			time.Sleep(10 * time.Second)
			continue
		}
		var msg Message

		err = json.Unmarshal([]byte(res), &msg)

		if err != nil {
			log.Println("MessageToMysql:", err)
			continue
		}

		SendAt, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Unix(msg.SendAt, 0).Format("2006-01-02 15:04:05"), time.Local)
		CreatedAt, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Unix(msg.CreatedAt, 0).Format("2006-01-02 15:04:05"), time.Local)

		tx := ay.Db.Begin()

		if tx.Table("im_msg").Create(&api.Msg{
			Type:      msg.Type,
			SendId:    msg.SendId,
			IsRead:    msg.IsRead,
			ReceiveId: msg.ReceiveId,
			Content:   msg.Content,
			SendAt:    models.MyTime{Time: SendAt},
			CreatedAt: models.MyTime{Time: CreatedAt},
		}).Error != nil {
			_, err = ay.Redis.Do("lpush", "message", res)
			tx.Rollback()
		} else {
			tx.Commit()
		}

	}

}
