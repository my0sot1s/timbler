package client

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	logx "github.com/my0sot1s/godef/log"
)

const timeOut = 20 // sec
var loop = 0

// Client is type CLient
type Client struct {
	connector *websocket.Conn
	clientURL string
	header    http.Header
}

// InitWsClient start client
func (c *Client) InitWsClient(clientURL string, header http.Header) {
	con, _, err := websocket.DefaultDialer.Dial(clientURL, header)
	if err != nil {
		logx.ErrLog(err)
	}
	c.connector = con
	c.clientURL = clientURL
	c.header = header
}

// Close is close connector
func (c *Client) Close() {
	c.connector.Close()
}

// Reconnect disconnector then create new
func (c *Client) Reconnect() {
	c.Close()
	c.InitWsClient(c.clientURL, c.header)
}

// ReadMsg read message
func (c *Client) ReadMsg(msgChan chan<- []byte) {
	done := make(chan struct{})
	defer close(done)
	for {
		_, message, err := c.connector.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			msgChan <- message
		}
		log.Printf("recv: %s", message)
	}
}

// SendMsg to server
func (c *Client) SendMsg(message []byte) error {
	err := c.connector.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		loop = loop + 1
		c.SendMsg(message)
	}
	loop = 0
	return err
}
