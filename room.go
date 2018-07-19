package main

import (
	"errors"

	"github.com/my0sot1s/tinker/utils"
)

// Room is a unit have many connection
type Room struct {
	Name    string
	ID      string
	Event   chan Event
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

}

func (r *Room) addClient(c *Connection) {
	for v := range r.Clients {
		if c.GetID() == v.GetID() {
			utils.Log("Connection is existed")
			return
		}
	}
	r.Clients[c] = true
}

func (r *Room) removeClient(c *Connection) {
	for v := range r.Clients {
		if c.GetID() == v.GetID() {
			c.connection.Close()
			delete(r.Clients, c)
			return
		}
	}
	utils.Log("Connection is existed")
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
			return
		}
	}
	utils.Log("Can not found room 1")
}
func (sr *StackRooms) removeMemberOfRoom(c *Connection, rname string) {
	for r := range sr.rooms {
		if r.GetName() == rname {
			r.removeClient(c)
			return
		}
	}
	utils.Log("Can not found room 2")
}

func (sr *StackRooms) checkRoomLen(rname string) int {
	for r := range sr.rooms {
		if r.GetName() == rname {
			return len(r.Clients)
		}
	}
	utils.Log("Can not found room 3")
	return -1
}

func (sr *StackRooms) addNewRoom(room *Room) {
	r := sr.rooms
	for R, state := range r {
		if R.ID == room.ID && state {
			utils.Log("Room is Existed")
			return
		}
	}
	sr.rooms[room] = true
	utils.Log("Room is Added: ", room.Name)
}

func (sr *StackRooms) removeRoom(rname string) {
	rooms := sr.rooms
	for r, state := range rooms {
		if r.GetName() == rname && state {
			delete(sr.rooms, r)
			utils.Log("Room is Deleted: ", r.Name)
			return
		}
	}
	utils.Log("Room is Not Existed: ", rname)
}
