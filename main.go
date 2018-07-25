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
	// StartCoreWs(port)
	rd := redislab.RedisCli{}
	rd.Config("redis-16703.c10.us-east-1-2.ec2.cloud.redislabs.com:16703", "redis-node", "95manhte")
	utils.Log("WS running :", port)
}
