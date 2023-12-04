/**
 * @Author: Hardews
 * @Date: 2023/4/5 22:01
 * @Description:
**/

package tool

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

var rDB *redis.Client

func RedisInit() {
	rDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rDB.Ping().Result()
	if err != nil {
		log.Println("redis service start failed")
	}
}

func Set(key, value string, expiration time.Duration) {
	rDB.Set(key, value, expiration)
}

func Get(key string) (string, error) {
	return rDB.Get(key).Result()
}

func Del(key string) {
	rDB.Del(key)
}
