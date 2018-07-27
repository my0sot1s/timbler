package redislab

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/my0sot1s/tinker/utils"
)

// RedisCli is struct of redis
type RedisCli struct {
	client   *redis.Client
	duration time.Duration
	mutex    *sync.Mutex
	message  chan *redis.Message
}
type callback func(string, string)

// Config init struct of redis
func (rc *RedisCli) Config(redisHost, redisDb, redisPass string) error {

	rc.client = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPass, // no password set
		DB:       0,         // use default DB
	})
	rc.duration = 60 * 10 * time.Second // 10 phut

	_, err := rc.client.Ping().Result()
	if err != nil {
		utils.ErrLog(err)
		return err
	}
	rc.mutex = &sync.Mutex{}
	rc.message = make(chan *redis.Message)
	go func() {
		for {
			mes := <-rc.message
			utils.Log("receive:", mes.Channel, " , ", mes.Payload)
		}
	}()
	utils.Log("ಠ‿ಠ Redis connected ಠ‿ಠ")
	return nil
}

// SetValue  to store
func (rc *RedisCli) SetValue(key string, value string, expiration time.Duration) error {
	return rc.client.Set(key, value, expiration).Err()
}

// GetValue from store
func (rc *RedisCli) GetValue(key string) (string, error) {
	val, err := rc.client.Get(key).Result()
	return val, err
}

// DelKey from store
func (rc *RedisCli) DelKey(key []string) (int, error) {
	val, err := rc.client.Del(key...).Result()
	return int(val), err
}

// LPushItem with key
func (rc *RedisCli) LPushItem(key string, timeExpired int, values ...interface{}) error {

	// str := make([]string, 0)
	for _, v := range values {
		b, e := json.Marshal(v)
		utils.ErrLog(e)
		_, err := rc.client.LPush(key, string(b)).Result()
		utils.ErrLog(err)
	}
	rc.SetExpired(key, timeExpired)
	return nil
}

// LRangeAll get with key
func (rc *RedisCli) LRangeAll(key string) ([]map[string]interface{}, error) {
	var raw []string
	raw, err := rc.client.LRange(key, 0, -1).Result()
	desMap := make([]map[string]interface{}, 0)
	for _, v := range raw {
		var d map[string]interface{}
		json.Unmarshal([]byte(v), &d)
		desMap = append(desMap, d)
	}
	return desMap, err
}

// SetExpired get time exprired key
func (rc *RedisCli) SetExpired(key string, min int) bool {
	d := time.Duration(min) * time.Minute

	b, err := rc.client.Expire(key, d).Result()
	utils.ErrLog(err)
	return b
}

// Publish message
func (rc *RedisCli) Publish(message string, channels []string) error {
	// rc.mutex.Lock()
	for _, channel := range channels {
		rc.client.Publish(channel, message)
	}
	// rc.mutex.Unlock()
	utils.Log("publish done")
	return nil
}

// Subscribe channels
func (rc *RedisCli) Subscribe(channels []string) *redis.PubSub {
	// rc.mutex.Lock()
	// defer rc.mutex.Unlock()
	p := rc.client.Subscribe(channels...)
	utils.Log("Subscribe done!")
	go rc.listenMessage(p)
	return p
}

func (rc *RedisCli) listenMessage(pubsub *redis.PubSub) {
	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}
		rc.message <- msg
	}
}

// Close redis
func (rc *RedisCli) Close() error {
	return rc.Close()
}