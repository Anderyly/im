package ay

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

var Redis redis.Conn

func ConnRedis() {

	c, err := redis.Dial("tcp", Yaml.GetString("redis.localhost")+":"+Yaml.GetString("redis.port"))
	if err != nil {
		log.Println("conn redis failed,", err)
		return
	}

	Redis = c
}
