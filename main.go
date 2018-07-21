package main

import (
	"os"

	"github.com/my0sot1s/tinker/utils"
)

var defaultPort = "8081"

func main() {
	utils.Log(utils.TitleConsole("T i m b l e r	 s t a r t e d"))
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	StartCoreWs(port)
	utils.Log("WS running :", port)
}
