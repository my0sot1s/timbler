package main

import "github.com/gin-gonic/gin"

// run ws
func runWs(rh *RoomHub, ctx *gin.Context) {
	c := &Connection{}
	c.InitConnection(rh, ctx.Writer, ctx.Request)
	go c.ReadMessageData()
	go c.WriteMessageData()
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
func StartCoreWs() {
	// Config Gin
	router := ginConfig()
	router.Static("/client", "./client")
	rh := &RoomHub{}
	go rh.Init()
	router.GET("/ws", func(ctx *gin.Context) {
		runWs(rh, ctx)
	})
	router.Run()
}
