package wsServer

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/my0sot1s/timbler/helper"
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
	// Messagetext
	Messagetext = "message"
)

// Message is struct message
type Message struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Created int    `json:"created"`
	// By is author
	By string `json:"by"`
	// To is room
	To string `json:"to"`
}

func (m Message) toByte() []byte {
	b, _ := json.Marshal(m)
	return b
}

// Event just message have name = `subscribe` || `unsubscribe`
type Event struct {
	Name    string `json:"name"`
	PayLoad string `json:"payload"`
}

// Connection is a user
type Connection struct {
	connection   *websocket.Conn
	ID           string
	send         chan []byte
	event        chan *Event
	mutex        *sync.Mutex
	roomsInvited map[*Room]bool
	roomHub      *RoomHub
}

// GetID just get ID
func (c Connection) GetID() string {
	return c.ID
}

// InitConnection create new Conn
func (c *Connection) InitConnection(rh *RoomHub, w http.ResponseWriter, r *http.Request) {
	// go c.listenEvent()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c.roomHub = rh
	c.connection = conn
	if c.ID == "" || &c.ID == nil {
		c.ID = helper.CreateID("cl")
	}
	c.mutex = &sync.Mutex{}
	c.roomsInvited = make(map[*Room]bool)
	c.send = make(chan []byte, 1024)
	c.event = make(chan *Event)
	c.sendId()
}

func (c *Connection) sendId() {
	m := Message{
		Type:    "system",
		Text:    c.GetID(),
		Created: time.Now().Second(),
		By:      "system",
		To:      "client",
	}
	message, _ := json.Marshal(m)
	c.send <- message
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

// IsConnectionInvited check connection is invited
func (c *Connection) IsConnectionInvited(rname string) bool {
	for r := range c.roomsInvited {
		if r.GetName() == rname {
			return true
		}
	}
	return false
}

// EventActived actived is process event
func (c *Connection) EventActived(evt *Event) {
	c.mutex.Lock()
	switch evt.Name {
	case Subscribe:
		// do something join
		// check room is Existed
		if isExisted := c.roomHub.IsRoomExisted(evt.PayLoad); !isExisted {
			// create room and add Room
			r := &Room{}
			r.createRoom(evt.PayLoad)
			c.roomHub.AddNewRoom(r)
		}
		r := c.roomHub.GetRoomByName(evt.PayLoad)
		if !c.IsConnectionInvited(r.GetName()) {
			c.roomsInvited[r] = true
		}
		r.Clients[c] = true

	case Unsubscribe:
		// do something leave
		if isExisted := c.roomHub.IsRoomExisted(evt.PayLoad); !isExisted {
			utils.Log("Room not existed false to leave")
		} else {
			r := c.roomHub.GetRoomByName(evt.PayLoad)
			c.roomsInvited[r] = false
			delete(c.roomsInvited, r)
			// check Room have agent
			if len := c.roomHub.ConnectionCountOnRoom(evt.PayLoad); len == 0 {
				c.roomHub.RemoveRoom(evt.PayLoad)
			}
			delete(r.Clients, c)
		}
	}
	c.mutex.Unlock()
}

// MessageFlowProcess split flow
func (c *Connection) MessageFlowProcess(bin []byte) {
	event := &Event{}
	utils.Str2T(string(bin), event)
	// utils.Log("Event Come == ", event.Name)
	switch event.Name {
	case Subscribe, Unsubscribe:
		slice := []string{}
		utils.Str2T(event.PayLoad, &slice)
		if event.Name == Subscribe {
			c.Subscribe(slice)
		} else if event.Name == Unsubscribe {
			c.Unsubscribe(slice)
		}
	case Messagetext:
		msg := &Message{}
		utils.Str2T(event.PayLoad, msg)
		utils.Log(msg, "blue")
		// 1. check room is existed room list
		r := c.roomHub.GetRoomByName(msg.To)
		if isSub := c.roomsInvited[r]; !isSub || r == nil {
			utils.Log("Room is not joined ", msg.To)
			return
		}
		// 2. send all client on room
		c.roomHub.SendMessageToRoom(r, msg)
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
			utils.ErrLog(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.CatchError(err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		utils.Log("message:", string(message))
		c.MessageFlowProcess(message)
		// do some thing with message
		// utils.Log("Add done!!")

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
		case evt, ok := <-c.event:
			if !ok {
				// c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			go c.EventActived(evt)
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
	utils.Log("On room ")
}

// Unsubscribe just definination any action leave rooms
func (c *Connection) Unsubscribe(rooms []string) {
	for _, r := range rooms {
		utils.Log("unsub: ", r)
		c.event <- &Event{
			Name:    Unsubscribe,
			PayLoad: r,
		}
	}
}
