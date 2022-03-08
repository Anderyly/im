package ay

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

var Redis redis.Conn

//
func init() {
	var yaml Yaml
	yaml.GetConf()
	c, err := redis.Dial("tcp", yaml.Redis.Localhost+":"+yaml.Redis.Port)
	if err != nil {
		log.Println("conn redis failed,", err)
		return
	}

	//defer c.Close()
	//_, err = c.Do("Set", "abc", 100)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	Redis = c
}
