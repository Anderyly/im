/*
 * *
 *  @Author anderyly
 *  @email admin@aaayun.cc
 *  @link http://blog.aaayun.cc/
 *  @copyright Copyright (c) 2022
 *  *
 */

package routers

import (
	"github.com/gin-gonic/gin"
	"im/controllers/ws"
)

func WsRouters(r *gin.RouterGroup) {

	wsGroup := r.Group("/ws")

	// 用户
	wsGroup.GET("", func(c *gin.Context) {
		ws.WsHandler(c.Writer, c.Request)
	})
	//wsGroup.GET("", ws.)

}
