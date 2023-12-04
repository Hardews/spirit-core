/**
 * @Author: Hardews
 * @Date: 2023/4/9 18:38
 * @Description:
**/

package service

type Hub struct {
	room       string
	num        int
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub(roomName string) *Hub {
	return &Hub{
		room:       roomName,
		num:        1,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.num--
			}
			if h.num == 0 {
				delete(Room, h.room)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					h.num--
				}
			}
		}
	}
}
