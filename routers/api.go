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
	"im/controllers/api"
)

func ApiRouters(r *gin.RouterGroup) {

	apiGroup := r.Group("/api/")

	// 用户
	apiGroup.POST("user/reg", api.UserController{}.Register)
	apiGroup.POST("user/login", api.UserController{}.Login)
	apiGroup.GET("user/D", api.UserController{}.D)
	apiGroup.POST("user/history", api.UserController{}.GetHistory)

}
