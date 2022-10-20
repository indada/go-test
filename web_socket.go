package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
	"time"
)

// gin 结合gorilla实现websocket

// // 解决跨域问题
var upGrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}
var clients = make(map[string]*websocket.Conn)

var msg = make(chan string, 1)

func client(c *gin.Context) {
	userId := c.Query("user")
	fmt.Println("client-userId", userId)
	//将id转为int64类型
	//gid, _ := strconv.ParseInt(groupId, 10, 64)
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	clients[userId] = ws
	defer ws.Close()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			delete(clients, userId)
			fmt.Println("ReadMessage", err)
			break
		}
		msg <- string(message)
		/*err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}*/
	}
}
func index(c *gin.Context) {
	t, _ := template.ParseFiles("template/index.html")
	tm := time.Now().Format("0405")
	dada := map[string]interface{}{
		"User": tm,
	}
	t.Execute(c.Writer, dada)
}

func main() {
	bindAddress := "localhost:2303"
	go Send()
	r := gin.Default()
	r.GET("client", client)
	r.GET("/", index)
	r.Run(bindAddress)

}
func Send() {
	println("Send")
	for {
		m := <-msg
		fmt.Println("接收到数据", m)
		for s, conn := range clients {
			err := conn.WriteJSON(m)
			fmt.Println("user:", s)
			if err != nil {
				fmt.Println("发送失败 user:", s, err)
			}
		}
	}
}
