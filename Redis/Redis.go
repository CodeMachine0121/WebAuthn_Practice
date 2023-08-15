package Redis

import (
	"Golang-WebAuthN/Utils"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
)

type Redis struct {
	Client *redis.Client
}

func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client
}

func (r *Redis) Set(key string, value interface{}) {

	jsonBytes, err := json.Marshal(value)
	jsonString := string(jsonBytes)
	Utils.ErrorHandle(err)

	err = r.Client.Set(key, jsonString, 0).Err() // => SET key value 0 數字代表過期秒數，在這裡0為永不過期
	Utils.ErrorHandle(err)
}

func (r *Redis) Get(key string, value interface{}) error {
	jsonString, err := r.Client.Get(key).Result()
	Utils.ErrorHandle(err)
	return json.Unmarshal([]byte(jsonString), value)
}
