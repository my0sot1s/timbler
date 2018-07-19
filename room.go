package main

import (
	"errors"

	"github.com/my0sot1s/tinker/utils"
)

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
	r.ID = createID("ro")
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
		c.send <- msg.toByte()
	}
}

// StackRooms is a lis room of connector
type StackRooms struct {
	rooms map[*Room]bool
}

func (sr *StackRooms) initStackRooms() {
	sr.rooms = make(map[*Room]bool)
}

func (sr *StackRooms) isExisted(name string) bool {
	for r := range sr.rooms {
		if r.GetName() == name {
			return true
		}
	}
	return false
}

func (sr *StackRooms) addMemberToRoom(c *Connection, rname string) {
	for r := range sr.rooms {
		if r.GetName() == rname {
			r.addClient(c)
			utils.Log("+ added client ", c.GetID(), " to ", r.GetName(), "green")
			return
		}
	}
	utils.Log("Can not found room: ", rname, "cyan")
}

func (sr *StackRooms) removeMemberOfRoom(c *Connection, rname string) {
	for r := range sr.rooms {
		if r.GetName() == rname {
			utils.Log("+ removed client ", c.GetID(), " to ", r.GetName(), "green")
			r.removeClient(c)
			return
		}
	}
	utils.Log("Can not found room ", rname, "cyan")
}

func (sr StackRooms) checkRoomLen(rname string) int {
	for r := range sr.rooms {
		if r.GetName() == rname {
			return len(r.Clients)
		}
	}
	utils.Log("Can not found room ", rname)
	return -1
}

func (sr *StackRooms) addNewRoom(room *Room) {
	r := sr.rooms
	for R, state := range r {
		if R.ID == room.ID && state {
			utils.Log("+ Room is Existed")
			return
		}
	}
	sr.rooms[room] = true
	utils.Log("+ Room is Added: ", room.Name)
}

func (sr *StackRooms) removeRoom(rname string) {
	rooms := sr.rooms
	for r, state := range rooms {
		if r.GetName() == rname && state {
			delete(sr.rooms, r)
			utils.Log("+ Room is Deleted: ", r.Name)
			return
		}
	}
	utils.Log("+ Room is Not Existed: ", rname)
}

func (sr *StackRooms) sendMessageToRoom(rname string, msg *Message) {
	for r := range sr.rooms {
		if r.GetName() == rname {
			go r.broadcast(msg)
		}
	}
	utils.Log("Not found room")
}
