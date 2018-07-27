package wsServer

import (
	"errors"
	"time"

	"github.com/my0sot1s/timbler/helper"

	"github.com/my0sot1s/tinker/utils"
)

// RoomHub room service
type RoomHub struct {
	rooms   map[*Room]bool
	created int
}

// Init roomHub
func (rh *RoomHub) Init() {
	rh.rooms = make(map[*Room]bool)
	rh.created = time.Now().Nanosecond()
}

// IsRoomExisted check room name is existed
func (rh RoomHub) IsRoomExisted(rname string) bool {
	for r := range rh.rooms {
		if r.GetName() == rname {
			return true
		}
	}
	return false
}

// GetRoomByName is get *Room by rname
func (rh *RoomHub) GetRoomByName(rname string) *Room {
	for r := range rh.rooms {
		if r.GetName() == rname {
			return r
		}
	}
	utils.Log("Room not existed")
	return nil
}

//ConnectionCountOnRoom count connection Online on room
func (rh RoomHub) ConnectionCountOnRoom(rname string) int {
	for r := range rh.rooms {
		if r.GetName() == rname {
			return len(r.Clients)
		}
	}
	utils.Log("Can not found room ", rname)
	return -1
}

// AddNewRoom is action add new Room
func (rh *RoomHub) AddNewRoom(room *Room) {
	for R := range rh.rooms {
		if R.GetName() == room.GetName() {
			utils.Log("+ Room is Existed")
			return
		}
	}
	rh.rooms[room] = true
	utils.Log("+ Room is Added: ", room.Name)
}

// RemoveRoom on hub
func (rh *RoomHub) RemoveRoom(rname string) {
	for r, state := range rh.rooms {
		if r.GetName() == rname && state {
			delete(rh.rooms, r)
			utils.Log("+ Room is Deleted: ", r.Name)
			return
		}
	}
	utils.Log("+ Room is Not Existed: ", rname)
}

//SendMessageToRoom is push message to client
func (rh *RoomHub) SendMessageToRoom(room *Room, msg *Message) {
	utils.Log("Sent to , ", room.GetName())
	for r := range rh.rooms {
		if r.GetName() == room.GetName() {
			utils.Log(len(r.Clients))
			go r.broadcast(msg)
			return
		}
	}
	utils.Log("Not found room")
}

// InjectEvent4Hub is push event to client
func (rh *RoomHub) InjectEvent4Hub(id string, event string, rooms []string) bool {
	for r := range rh.rooms {
		if isDone := r.find4SubOrUnsub(id, event, rooms); isDone {
			return isDone
		}
	}
	return false
}

// Room is a unit have many connection
type Room struct {
	Name    string
	ID      string
	Clients map[*Connection]bool
}

// GetID is get room ID
func (r Room) GetID() string {
	return r.ID
}

// GetName is get room Name
func (r Room) GetName() string {
	return r.Name
}
func (r *Room) createRoom(name string) {
	if name == "" {
		utils.ErrLog(errors.New("No room name"))
		return
	}
	r.Name = name
	r.ID = helper.CreateID("ro")
	r.Clients = make(map[*Connection]bool)

}

func (r *Room) addClient(c *Connection) {
	for v := range r.Clients {
		if c.GetID() == v.GetID() {
			utils.Log("++ Connection is existed", "cyan")
			return
		}
	}
	utils.Log("++ Add client success ", c.GetID(), "green")
	r.Clients[c] = true
}

func (r *Room) removeClient(c *Connection) {
	for v := range r.Clients {
		if c.GetID() == v.GetID() {
			// c.connection.Close()
			utils.Log("++ Deleted client success ", c.GetID(), "green")
			delete(r.Clients, c)
			return
		}
	}
	utils.Log("Connection is existed", "cyan")
}

func (r *Room) broadcast(msg *Message) {
	for c := range r.Clients {
		utils.Log("--->", string(msg.toByte()))
		c.send <- msg.toByte()
	}
}

func (r *Room) find4SubOrUnsub(id string, event string, rooms []string) bool {
	for c := range r.Clients {
		if c.GetID() != id {
			continue
		}
		if event == "subscribe" {
			c.Subscribe(rooms)
		} else if event == "unsubscribe" {
			c.Unsubscribe(rooms)
		}
		return true
	}
	return false
}
