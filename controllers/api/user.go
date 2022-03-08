/*
 * *
 *  @Author anderyly
 *  @email admin@aaayun.cc
 *  @link http://blog.aaayun.cc/
 *  @copyright Copyright (c) 2022
 *  *
 */

package api

import (
	"github.com/gin-gonic/gin"
	"im/ay"
	"im/models/api"
)

type UserController struct {
}

type GetRegisterForm struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
	Nickname string `form:"nickname" binding:"required"`
}

// Register 注册
func (con UserController) Register(c *gin.Context) {

	var data GetRegisterForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, "400", ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var user api.User
	ay.Db.First(&user, "account = ?", data.Account)

	if user.ID == 0 {
		res := ay.Db.Create(&api.User{
			Account:  data.Account,
			Password: ay.MD5(data.Password),
			NickName: data.Nickname,
			Status:   1,
		})
		if res.Error != nil {
			ay.Json{}.Msg(c, "400", "创建失败", gin.H{
				"error": res.Error,
			})
		} else {
			ay.Json{}.Msg(c, "200", "创建成功", gin.H{
				"token": ay.AuthCode(data.Account, "ENCODE", "", 0),
			})
		}

	} else {
		ay.Json{}.Msg(c, "1001", "账号已存在", gin.H{})
	}

}

type GetLoginForm struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Login 登入
func (con UserController) Login(c *gin.Context) {
	var data GetLoginForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, "400", ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var user api.User
	ay.Db.First(&user, "account = ?", data.Account)

	if user.ID == 0 {
		ay.Json{}.Msg(c, "400", "账号不存在", gin.H{})
		return
	}

	if user.Password != ay.MD5(data.Password) {
		ay.Json{}.Msg(c, "400", "密码错误", gin.H{})
		return
	}

	if user.Status == 0 {
		ay.Json{}.Msg(c, "400", "账号已被冻结", gin.H{})
		return
	}

	ay.Json{}.Msg(c, "200", "success", gin.H{
		"token": ay.AuthCode(user.Account, "ENCODE", "", 0),
	})
}

type GetDForm struct {
	Str string `form:"str" binding:"required"`
}

// Login 登入
func (con UserController) D(c *gin.Context) {
	var data GetDForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, "400", ay.Validator{}.Translate(err), gin.H{})
		return
	}

	res := ay.AuthCode(data.Str, "DECODE", "", 0)

	ay.Json{}.Msg(c, "200", "success", gin.H{
		"token": res,
	})
}

type GetHistoryForm struct {
	Token     string `form:"token" binding:"required"`
	ReceiveId string `form:"receive_id" binding:"required"`
	Page      int    `form:"page" binding:"required"`
}

// GetHistory 获取历史
func (con UserController) GetHistory(c *gin.Context) {
	var data GetHistoryForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, "400", ay.Validator{}.Translate(err), gin.H{})
		return
	}

	account := ay.AuthCode(data.Token, "DECODE", "", 0)

	if account == "" {
		ay.Json{}.Msg(c, "400", "token错误", gin.H{})
		return
	}

	var user api.User
	ay.Db.First(&user, "account = ?", account)

	if user.ID == 0 {
		ay.Json{}.Msg(c, "400", "账号不存在", gin.H{})
		return
	}

	if user.Status == 0 {
		ay.Json{}.Msg(c, "400", "账号已被冻结", gin.H{})
		return
	}

	page := data.Page - 1

	var msg []api.Msg
	ay.Db.Limit(20).Offset(page*20).Where("type = 1 AND (send_id = ? AND receive_id = ?) or (send_id = ? AND receive_id = ?)", account, data.ReceiveId, data.ReceiveId, account).Order("created_at desc").Find(&msg)

	ay.Json{}.Msg(c, "200", "success", gin.H{
		"list": msg,
	})
}
