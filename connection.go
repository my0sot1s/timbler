package main

import (
	"bytes"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/my0sot1s/tinker/utils"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  maxMessageReadSize * 2,
		WriteBufferSize: maxMessageReadSize * 2,
	}
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	maxMessageReadSize = 1024
	pongWait           = 60 * time.Second
	writeWait          = 10 * time.Second
	pingPeriod         = (pongWait * 9) / 10
	// Subscribe join name
	Subscribe = "subscribe"
	// Unsubscribe leave name
	Unsubscribe = "unsubscribe"
	// MessageText
	MessageText = "message"
)

// Message is struct message
type Message struct {
	Text    string `json:"text"`
	Created int    `json:"created"`
	// By is author
	By string `json:"by"`
	// To is room
	To string `json:"to"`
}

// Event just message have name = `subscribe` || `unsubscribe`
type Event struct {
	Name    string `json:"name"`
	PayLoad string `json:"payload"`
}

// Connection is a user
type Connection struct {
	connection *websocket.Conn
	ID         string
	send       chan []byte
	event      chan *Event
	mutex      *sync.Mutex
	rooms      *StackRooms
}

// GetID just get ID
func (c Connection) GetID() string {
	return c.ID
}

// InitConnection create new Conn
func (c *Connection) InitConnection(w http.ResponseWriter, r *http.Request) {
	go c.listenEvent()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c.connection = conn
	if c.ID == "" || &c.ID == nil {
		c.ID = createID("cl")
	}
	c.mutex = &sync.Mutex{}
	s := &StackRooms{}
	s.initStackRooms()
	c.rooms = s

}

func (c *Connection) listenEvent() {
	for {
		select {
		case event := <-c.event:
			utils.Log(event)
			c.EventActived(event)
		}
	}
}

// SetReadDeadline  just Add wait
func (c *Connection) SetReadDeadline() {
	c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

// SetWriteDeadline just Add wait
func (c *Connection) SetWriteDeadline() {
	c.connection.SetWriteDeadline(time.Now().Add(writeWait))
}

// CatchError is action Log Error
func (c *Connection) CatchError(err error) {
	utils.ErrLog(err)
}

// EventActived actived is process event
func (c *Connection) EventActived(evt *Event) {
	switch evt.Name {
	case Subscribe:
		c.mutex.Lock()
		// do something join
		// check room is Existed
		if isExisted := c.rooms.isExisted(evt.Name); !isExisted {
			// create room and add Room
			r := &Room{}
			r.createRoom(evt.Name)
			c.rooms.addNewRoom(r)
		}
		// add member to Room
		c.rooms.addMemberToRoom(c, evt.Name)
		c.mutex.Unlock()
		break
	case Unsubscribe:
		c.mutex.Lock()
		// do something leave
		if isExisted := c.rooms.isExisted(evt.Name); !isExisted {
			utils.Log("Room not existed false to leave")
			return
		}
		c.rooms.removeMemberOfRoom(c, evt.Name)
		// check Room have agent
		if len := c.rooms.checkRoomLen(evt.Name); len == 0 {
			c.rooms.removeRoom(evt.Name)
		}
		c.mutex.Unlock()
		break
	}
}

// MessageFlowProcess split flow
func (c *Connection) MessageFlowProcess(bin []byte) {
	event := &Event{}
	utils.Str2T(string(bin), event)
	utils.Log("Event Come == ", event.Name)
	switch event.Name {
	case Subscribe, Unsubscribe:
		slice := []string{}
		utils.Str2T(event.PayLoad, &slice)
		if event.Name == Subscribe {
			c.Subscribe(slice)
		} else if event.Name == Unsubscribe {
			c.Unsubscribe(slice)
		}
	case MessageText:
		msg := &Message{}
		utils.Str2T(event.PayLoad, msg)
	default:
		utils.Log("come fuck")
	}

}

// ReadMessageData is a action Read message from any room
func (c *Connection) ReadMessageData() {
	defer func() {
		c.connection.Close()
	}()
	c.connection.SetReadLimit(maxMessageReadSize)
	c.SetReadDeadline()
	c.connection.SetPongHandler(func(string) error {
		c.SetReadDeadline()
		return nil
	})
	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.CatchError(err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		utils.Log("message:", string(message))
		c.MessageFlowProcess(message)
		// do some thing with message
	}
}

// WriteMessageData is broad cast data to a room
func (c *Connection) WriteMessageData() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.SetWriteDeadline()
			if !ok {
				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				c.CatchError(err)
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		// case evt, ok := <-c.event:
		// 	if !ok {
		// 		// c.connection.WriteMessage(websocket.CloseMessage, []byte{})
		// 		return
		// 	}
		// 	c.EventActived(evt)
		case <-ticker.C:
			c.SetWriteDeadline()
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Subscribe is a action definition of join any rooms
func (c *Connection) Subscribe(rooms []string) {
	for _, r := range rooms {
		utils.Log("sub: ", r)
		c.event <- &Event{
			Name:    Subscribe,
			PayLoad: r,
		}
	}
}

// Unsubscribe just definination any action leave rooms
func (c *Connection) Unsubscribe(rooms []string) {
	for _, r := range rooms {
		c.event <- &Event{
			Name:    Unsubscribe,
			PayLoad: r,
		}
	}
}
