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
	"gorm.io/gorm"
)

//type UserModel struct {
//}

type User struct {
	gorm.Model
	Account  string `gorm:"column:account" json:"account"`
	Password string `gorm:"column:password" json:"password"`
	NickName string `gorm:"column:nickname" json:"nickname"`
	Status   int    `gorm:"column:status" json:"status"`
}

func (User) TableName() string {
	return "im_users"
}
