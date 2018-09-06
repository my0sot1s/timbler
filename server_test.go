package timbler

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestServer(t *testing.T) {
	// start server
	serve := gin.Default()
	ws := &RealtimeWS{}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	ws.InitWS(serve, dir)

	serve.Run(":8080")

}
