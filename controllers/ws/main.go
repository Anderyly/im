package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"im/ay"
	"log"
	"net/http"
	"strconv"
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
}

var clientMap map[string]*Node = make(map[string]*Node, 0)

// Message {"send_id":"123","type":1,"receive_id":"456","content":""}
type Message struct {
	Id        int64  `json:"id,omitempty" form:"id"`                 // 消息ID
	Type      int    `json:"type,omitempty" form:"type"`             // 消息类型： 1私聊 9检测心跳
	SendId    string `json:"send_id,omitempty" form:"send_id"`       // 发送者
	ReceiveId string `json:"receive_id,omitempty" form:"receive_id"` // 接收者
	Content   string `json:"content,omitempty" form:"content"`       // 消息的内容
	CreatedAt int64  `json:"created_at,omitempty" form:"created_at"` // 消息的内容
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		conn *websocket.Conn
		err  error
	)
	query := r.URL.Query()
	token := query.Get("token")
	log.Println(token)
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
	}

	clientMap[id] = node

	go Send(node)
	go Receive(node)

	msg := `{"send_id":"System","type":0,"receive_id":"` + id + `","content":"用户：` + id + ` 接入成功", "created_at":` + strconv.FormatInt(time.Now().Unix(), 10) + `}`

	// 设置链接用户在线
	UserOnline(id)

	SendMsg(id, []byte(msg))
}

func Dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}

	errMsg := Message{
		Type:      0,
		SendId:    "System",
		ReceiveId: msg.SendId,
		Content:   "发送人不可与接收人一致",
		CreatedAt: time.Now().Unix(),
	}

	// 发送人处理
	if msg.ReceiveId == msg.SendId {
		output, _ := json.Marshal(errMsg)
		SendMsg(msg.SendId, output)
	}

	switch msg.Type {
	case 1: // 私聊
		Chat{}.Operate(msg, data)
	case 9: // 检测客户端的心跳
		// 更新用户在线状态
		UserOnline(msg.SendId)
	}
}
