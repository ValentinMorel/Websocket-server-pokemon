package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	var base_addr = "localhost:8081"
	//var flag int
	u := url.URL{Scheme: "ws", Host: base_addr, Path: "/socket"}
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
	defer c.Close()
	quit := make(chan bool)
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			/*
				if flag == 0 {
					time.Sleep(3 * time.Second)
					s := fmt.Sprintf(`{"event":"update", "id": "0909090909", "mapId": "33", "x": 5, "y": 1, "spriteId": 3}`)
					c.WriteMessage(websocket.TextMessage, []byte(s))
					flag += 1
				}
			*/
			log.Println("recv:", string(message))

		}
	}()
	<-quit

}
