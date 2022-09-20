package ws

import (
	"github.com/gorilla/websocket"
	"im/ay"
	"log"
)

func Send(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("Send:" + err.Error())
				goto ERR
				//return
			}
		}

	}
ERR:
	ay.Redis.Do("ZREM", "online_users", node.Id)
	node.Conn.Close()

}

func SendMsg(userId string, msg []byte) {
	node, ok := clientMap[userId]
	if ok {
		node.DataQueue <- msg
	}
}
