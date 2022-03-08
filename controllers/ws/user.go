package ws

import (
	"im/ay"
	"log"
	"time"
)

func UserOnline(id string) {
	_, err := ay.Redis.Do("ZADD", "online_users", id, time.Now().Unix())
	if err != nil {
		log.Println(err)
		return
	}
}
