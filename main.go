package main

import (
	"github.com/gin-gonic/gin"
	"im/ay"
	"im/controllers/ws"
	"im/routers"
	"im/service"
	"net/http"
)

var r *gin.Engine

func main() {

	//go ws.Manager.Start()

	ay.Yaml = ay.InitConfig()
	ay.Sql()
	ay.ConnRedis()
	go ay.WatchConf()

	r = gin.Default()
	r = service.Set(r)
	r.StaticFS("/static/", http.Dir("./static"))
	r = routers.GinRouter(r)

	// 开启redis查询
	go ws.MessageToMysql()

	err := r.Run(":" + ay.Yaml.GetString("server.port"))
	if err != nil {
		panic(err.Error())
	}
}
