package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"im/ay"
	"log"
	"net/http"
	"time"
)

var (
	upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	//GroupSets set.Interface
	Id string
}

var clientMap map[string]*Node = make(map[string]*Node, 0)

type Message struct {
	Id        int64  `json:"id,omitempty" form:"id"`                 // 消息ID
	Type      int    `json:"type,omitempty" form:"type"`             // 消息类型： 1私聊 9检测心跳
	IsRead    int    `json:"is_read,omitempty" form:"id_read"`       // 消息是否已读
	SendId    string `json:"send_id,omitempty" form:"send_id"`       // 发送者
	ReceiveId string `json:"receive_id,omitempty" form:"receive_id"` // 接收者
	Content   string `json:"content,omitempty" form:"content"`       // 消息的内容
	SendAt    int64  `json:"send_at"`                                // 时间戳
	CreatedAt int64  `json:"created_at,omitempty" form:"created_at"` // 时间戳
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		conn *websocket.Conn
		err  error
	)
	query := r.URL.Query()
	token := query.Get("token")

	id := ay.AuthCode(token, "DECODE", "", 0)
	if id == "" {
		log.Println("用户token错误")
		return
	}

	if conn, err = upgrade.Upgrade(w, r, nil); err != nil {
		log.Printf("ws链接错误：%s", err)
		return
	}

	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		Id:        id,
	}

	clientMap[id] = node

	go Send(node)
	go Receive(node)

	UserOnline(id)
}

func Dispatch(id string, data []byte) {

	log.Println(id)
	UserOnline(id)

	if string(data) == "ping" {
		SendMsg(id, []byte("pong"))
		return
	}

	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println("Dispatch:" + err.Error())
		return
	}

	errMsg := Message{
		Type:      0,
		SendId:    "System",
		ReceiveId: msg.SendId,
		CreatedAt: time.Now().Unix(),
	}

	// 发送人处理
	if id != msg.SendId {
		errMsg.Content = "请检查send_id"
		output, _ := json.Marshal(errMsg)
		SendMsg(id, output)
		return
	}

	// 发送人处理
	if msg.ReceiveId == msg.SendId {
		errMsg.Content = "发送人不可与接收人一致"
		output, _ := json.Marshal(errMsg)
		SendMsg(id, output)
		return
	}

	// 空消息处理
	if msg.ReceiveId == msg.SendId {
		errMsg.Content = "发送消息不能为空"
		output, _ := json.Marshal(errMsg)
		SendMsg(id, output)
		return
	}

	switch msg.Type {
	case 1: // 私聊
		Chat{}.Operate(msg, data)
	}
}
