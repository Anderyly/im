package ws

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"im/ay"
	"log"
	"time"
)

func Receive(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		Dispatch(data)
		fmt.Printf("recv<=%s", data)
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
			log.Println("unmarshal message failed,", err)
			continue
		}

		ay.Db.Table("im_msg").Create(&msg)
		//log.Println(msg)
	}

}
