package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
	"time"
)
var (
	// 定义常量
	RedisClient     *redis.Pool
	REDIS_HOST string
	REDIS_DB   int
)

func TestLinkRedis(t *testing.T){
	REDIS_HOST = "localhost:6379"
	REDIS_DB = 0
	RedisClient = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}

	rc := RedisClient.Get()
	rule :="rule1"
	device :="8"
	whiteL := "1,8,9,y,i,k,3,r,44,kk"

	luaPull := redis.NewScript(3,CACHE_AND_SEARCH)
	//status,_ :=redis.Int(luaPull.Do(rc,rule,device,",",whiteList))
	status,_ := redis.Int( luaPull.Do(rc,rule,device,",",whiteL))
	fmt.Println(status)

	defer rc.Close()
}