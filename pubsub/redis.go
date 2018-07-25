package redis

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/my0sot1s/tinker/utils"
)

type RedisCli struct {
	client   *redis.Client
	duration time.Duration
	mutex    *sync.Mutex
	message  chan []byte
}

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
	rc.message = make(chan []byte)
	utils.Log("ಠ‿ಠ Redis connected ಠ‿ಠ")
	go rc.listenMessage()
	return nil
}
func (rc *RedisCli) SetValue(key string, value string, expiration time.Duration) error {
	return rc.client.Set(key, value, expiration).Err()
}

func (rc *RedisCli) GetValue(key string) (string, error) {
	val, err := rc.client.Get(key).Result()
	return val, err
}

func (rc *RedisCli) DelKey(key []string) (int, error) {
	val, err := rc.client.Del(key...).Result()
	return int(val), err
}

// multiple item

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

func (rc *RedisCli) SetExpired(key string, min int) bool {
	d := time.Duration(min) * time.Minute

	b, err := rc.client.Expire(key, d).Result()
	utils.ErrLog(err)
	return b
}

func (rc *RedisCli) Publish(message string, channels []string) error {
	for _, channel := range channels {
		rc.mutex.Lock()
		rc.client.Publish(channel, message)
		rc.mutex.Unlock()
	}
	return nil
}

func (rc *RedisCli) Subscribe(channels []string) error {
	rc.mutex.Lock()
	rc.client.Subscribe(channels...)
	defer rc.mutex.Unlock()
	return nil
}

func (rc *RedisCli) listenMessage() {
	for {
	}
}

func (rc *RedisCli) Close() error {
	return rc.Close()
}
