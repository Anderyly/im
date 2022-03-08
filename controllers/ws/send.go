package ws

import (
	"github.com/gorilla/websocket"
	"log"
)

func Send(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

func SendMsg(userId string, msg []byte) {
	node, ok := clientMap[userId]
	if ok {
		node.DataQueue <- msg
	}
}
