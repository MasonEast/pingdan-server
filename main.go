package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Orders struct {
	Price string `json:"price"`
	Type string `json:"type"`
	User string `json:"user"`
	RealPrice int `json:"realPrice"`
}
type Message struct {
	Message []Orders `json:"message"`
}

var upgrader = websocket.Upgrader{
	//check origin will check the cross region source (note : please not using in production)
	CheckOrigin: func(r *http.Request) bool {
		//Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
		return true
	},
}

func main() {
	r := gin.Default()
	hub := NewHub()

	go hub.run()

	r.GET("/ws", func(c *gin.Context) {
		//upgrade get request to websocket protocol
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func() {
			delete(hub.clients, ws)
			ws.Close()
			fmt.Printf("closed")
		}()

		hub.clients[ws] = true
		fmt.Println("connected")
		read(hub, ws)

		// for {
			//Read Message from client
			// mt, message, err := ws.ReadMessage()

			//If client message is ping will return pong
			// if string(message) == "ping" {
			// 	message = []byte("pong")
			// }
			// //Response message to client
			// err = ws.WriteMessage(mt, message)
			// if err != nil {
			// 	fmt.Println(err)
			// 	break
			// }
		// }
	})
	r.Run(":5000") // listen and serve on 0.0.0.0:8080
}

func read(hub *Hub, client *websocket.Conn) {
	for {
		var message Message
		
		err := client.ReadJSON(&message)
		fmt.Println(message)
		if err != nil {
			delete(hub.clients, client)
			fmt.Printf("errror occurred: %v:", err)
			break
		}
		hub.broadcast <- message
	}
}
