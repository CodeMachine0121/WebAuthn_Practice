package Redis

import (
	"Golang-WebAuthN/Models"
	"Golang-WebAuthN/Utils"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
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

func (r *Redis) Set(key string, value Models.User) {

	jsonBytes, err := json.Marshal(value)
	jsonString := string(jsonBytes)
	Utils.ErrorHandle(err)

	err = r.Client.Set(key, jsonString, 0).Err() // => SET key value 0 數字代表過期秒數，在這裡0為永不過期
	Utils.ErrorHandle(err)
}

func (r *Redis) Get(key string) string {
	jsonString, err := r.Client.Get(key).Result()
	Utils.ErrorHandle(err)
	return jsonString
}

func (r *Redis) JsonParseUser(jsonString string) Models.User {

	var user Models.User
	if err := json.Unmarshal([]byte(jsonString), &user); err != nil {
		log.Fatal("Parsing error")
	}
	return user
}
