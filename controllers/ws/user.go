package ws

import (
	"im/ay"
	"log"
	"time"
)

func UserOnline(id string) {
	ay.Redis.Do("ZREM", "online_users", id)
	_, err := ay.Redis.Do("ZADD", "online_users", time.Now().Unix(), id)
	if err != nil {
		log.Println(err)
		return
	}
}
