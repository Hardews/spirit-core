/**
 * @Author: Hardews
 * @Date: 2023/4/9 18:38
 * @Description:
**/

package service

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"spirit-core/my_consts"
	"spirit-core/tool"
	"time"
)

var (
	newline         = []byte{'\n'}
	space           = []byte{' '}
	ErrOfNoThisRoom = errors.New("no this room")
)

var Room = make(map[string]*Hub)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	roomName string
	username string
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
}

func NewRoom(ctx *gin.Context, username string) (string, error) {
	roomName := RandomStr(9)

	hub := newHub(roomName)
	go hub.run()

	Room[roomName] = hub

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return "", err
	}

	client := &Client{
		roomName: roomName,
		username: username,
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 1024),
	}

	hub.register <- client

	client.send <- []byte(client.username + " 创建了房间,房间号为:" + client.roomName)
	go client.readPump()
	go client.writePump()

	return roomName, nil
}

func JoinRoom(ctx *gin.Context, roomName, username string) error {
	hub, ok := Room[roomName]
	if !ok {
		return ErrOfNoThisRoom
	}
	go hub.run()
	hub.num++

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("join room error, room name:%s, err:%s", roomName, err.Error())
		return err
	}

	client := &Client{
		username: username,
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 1024),
	}

	hub.register <- client

	client.send <- []byte(client.username + "加入了房间")
	go client.readPump()
	go client.writePump()

	return nil
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(my_consts.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(my_consts.PongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(my_consts.PongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		result, err := tool.ContentAudit(string(message))
		if err != nil {
			log.Printf("content check failed, err:%s, room name:%s", err.Error(), c.roomName)
		}

		if result.ConclusionType != 1 {
			message = []byte("**********")
		}

		message = bytes.TrimSpace(bytes.Replace([]byte(c.username+":"+string(message)), newline, space, -1))
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(my_consts.PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(my_consts.WriteWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(my_consts.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWSYZ1234567890abcdefghijklmnopqrstuvwsyz1234567890")

// RandomStr 生成随机字符串
func RandomStr(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
