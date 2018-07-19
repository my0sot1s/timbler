package main

import "github.com/gin-gonic/gin"

// run ws
func runWs(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	c := &Connection{}
	c.InitConnection(w, r)
	go c.ReadMessageData()
	go c.WriteMessageData()
}

func ginConfig() *gin.Engine {

	mode := gin.TestMode
	// set mode `production` or `dev`
	gin.SetMode(mode)
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(gin.Logger())
	return g
}

// StartCoreWs Start Ws
func StartCoreWs() {
	// Config Gin
	router := ginConfig()
	router.Static("/client", "./client")
	router.GET("/ws", runWs)
	router.Run()
}
