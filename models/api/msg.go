/*
 * *
 *  @Author anderyly
 *  @email admin@aaayun.cc
 *  @link http://blog.aaayun.cc/
 *  @copyright Copyright (c) 2022
 *  *
 */

package api

import "im/models"

//type UserModel struct {
//}

type Msg struct {
	Id        int64         `json:"-"`
	Type      int           `gorm:"column:type" json:"type"`
	SendId    string        `gorm:"column:send_id" json:"send_id"`
	IsRead    int           `gorm:"column:is_read" json:"is_read"`
	ReceiveId string        `gorm:"column:receive_id" json:"receive_id"`
	Content   string        `gorm:"column:content" json:"content"`
	SendAt    models.MyTime `gorm:"column:send_at" json:"send_at"`
	CreatedAt models.MyTime `gorm:"column:created_at" json:"created_at"`
}

func (Msg) TableName() string {
	return "im_msg"
}
