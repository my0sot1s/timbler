package main

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/my0sot1s/tinker/utils"
)

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
	if isDone := hub.injectEvent4Hub(id, event, rooms); !isDone {
		utils.Log("Can not find connection by cid:", id)
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

func ginConfig() *gin.Engine {

	mode := gin.TestMode
	// set mode `production` or `dev`
	gin.SetMode(mode)
	g := gin.New()
	g.Use(gin.Recovery(), gin.Logger())
	return g
}

// StartCoreWs Start Ws
func StartCoreWs(port string) {
	// Config Gin
	router := ginConfig()
	router.Static("/client", "./client")
	rh := &RoomHub{}
	go rh.Init()
	router.GET("/ws", func(ctx *gin.Context) {
		runWs(rh, ctx)
	})
	router.POST("/ws/sub-unsub", func(ctx *gin.Context) {
		subUnSub(rh, ctx)
	})

	router.Run(":" + port)
}
