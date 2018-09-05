package timbler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	logx "github.com/my0sot1s/godef/log"
)

// RealtimeWS realtime struct
type RealtimeWS struct {
	rh *RoomHub
}

// run ws
func runWs(rh *RoomHub, ctx *gin.Context) {
	c := &Connection{}
	c.InitConnection(rh, ctx.Writer, ctx.Request)
	go c.ReadMessageData()
	go c.WriteMessageData()
}

// for data = {id:string,event:string,payload:[{name:string,payload:string}]}
func subUnSub(hub *RoomHub, ctx *gin.Context) {
	id := ctx.PostForm("id")
	event := ctx.PostForm("event")
	payload := ctx.PostForm("payload")
	var rooms []string
	json.Unmarshal([]byte(payload), &rooms)
	if event != "subscribe" && event != "unsubscribe" {
		ctx.JSON(400, gin.H{
			"name":    "error",
			"payload": "Event is not valid",
		})
		return
	}
	if isDone := hub.InjectEvent4Hub(id, event, rooms); !isDone {
		logx.Log("Can not find connection by cid:", id)
		ctx.JSON(400, gin.H{
			"name":    "error",
			"payload": "Can not find connection by cid:",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"name": "success",
	})
}

// InitWS create form
func (ws *RealtimeWS) InitWS(router *gin.Engine, root string) {
	ws.rh = &RoomHub{}
	go ws.rh.Init()

	logx.Log("-- WebSocket started --:")

	router.Static("/client", root+"/client")

	router.GET("/ws", func(ctx *gin.Context) {
		runWs(ws.rh, ctx)
	})

	router.POST("/ws/sub-unsub", func(ctx *gin.Context) {
		subUnSub(ws.rh, ctx)
	})
}
