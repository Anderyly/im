package main

import (
	"github.com/gin-gonic/gin"
	"im/controllers/ws"
	"im/routers"
	"im/service"
	"net/http"
)

var r *gin.Engine

func main() {

	//go ws.Manager.Start()

	r = gin.Default()

	r = service.Set(r)

	//r.LoadHTMLGlob("views/**/**/*")
	r.StaticFS("/static/", http.Dir("./static"))

	r = routers.GinRouter(r)
	//ay.Con()
	go ws.MessageToMysql()
	//http.HandleFunc("/ws", ws.MainController{}.Home)
	err := r.Run(":8090")
	if err != nil {
		panic(err.Error())
	}
}
