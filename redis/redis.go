package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
)

//Connection stores the current redis connection
var connection *redis.Client

//Init creates the redis client
func Init() {
	dev.LogInfo("Redis is being initialized")
	connection = redis.NewClient(&redis.Options{
		Addr: fmt.Sprint(config.CurrentConfig.Redis.Host, ":", config.CurrentConfig.Redis.Port),
		DB:   config.CurrentConfig.Redis.Database, // use default DB
	})

	_, err := connection.Ping().Result()
	if err != nil {
		dev.LogFatal("Couldn't connect to redis: ", err.Error())
	}

	dev.LogInfo("Connection to redis established")
}

//SessionValid return true if the session key is found
func SessionValid(Session string) bool {
	key, _ := connection.Get(Session).Result()
	if len(key) == 0 {
		return false
	}

	go refreshSession(Session)
	return true
}

func refreshSession(Session string) {
	//Value is just the user id
	//The user id is contained inside the sessionkey (Example: session_1_as231fsdf413 -> 1 will be the user id)
	var Value = strings.Split(strings.Split(Session, "session_")[1], "_")[0]
	if err := connection.Set(Session, Value, time.Hour).Err(); err != nil {
		dev.LogError("Couldn't refresh session \""+Session+"\" -> ", err.Error())
	}
}
