package main

import (
	"fmt"
	"os"
	"time"

	"github.com/my0sot1s/timbler/redislab"
	"github.com/my0sot1s/tinker/utils"
)

var defaultPort = "8081"

func runRedis() {
	rd := redislab.RedisCli{}
	c := make(chan bool, 1)
	go func(cx chan bool) {
		rd.Config("redis-16703.c10.us-east-1-2.ec2.cloud.redislabs.com:16703", "redis-node", "95manhte")
		rd.Subscribe([]string{"chan1"})
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			rd.Publish(fmt.Sprintf("Hello %d", i), []string{"chan1"})
		}
	}(c)
	<-c
}
func main() {
	utils.Log(utils.TitleConsole("T i m b l e r	 s t a r t e d"))
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	StartCoreWs(port)
	// runRedis()
	utils.Log("WS running :", port)
}
