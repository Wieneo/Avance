package redis

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
)

//Connection stores the current redis connection
var connection *redis.Client

//Init creates the redis client
func Init() {
	dev.LogInfo("Redis is being initialized")

	if !config.CurrentConfig.Redis.Sentinel.Enabled {
		connection = redis.NewClient(&redis.Options{
			Addr: fmt.Sprint(config.CurrentConfig.Redis.Host, ":", config.CurrentConfig.Redis.Port),
			DB:   config.CurrentConfig.Redis.Database, // use default DB
		})
	} else {
		connection = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    config.CurrentConfig.Redis.Sentinel.Master,
			SentinelAddrs: config.CurrentConfig.Redis.Sentinel.Endpoints,
			DB:            config.CurrentConfig.Redis.Database,
		})
	}

	_, err := connection.Ping().Result()
	if err != nil {
		dev.LogFatal(err, "Couldn't connect to redis: ", err.Error())
	}

	dev.LogInfo("Connection to redis established")
}

//Ping returns true if redis is alive
func Ping() (bool, error) {
	_, err := connection.Ping().Result()
	if err != nil {
		return false, err
	}
	return true, nil
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

//CreateSession returns a session key for the given userid
//SessionKey will be saved as follows: session_USERID_UUID
//This is used to be able to delete all sessions of a user
func CreateSession(UserID int) (string, error) {
	Session := uuid.New()
	if err := connection.Set("session_"+strconv.Itoa(UserID)+"_"+Session.String(), strconv.Itoa(UserID), time.Hour).Err(); err != nil {
		return "", err
	}
	return "session_" + strconv.Itoa(UserID) + "_" + Session.String(), nil
}

//DestroySession removes the session key from redis
func DestroySession(SessionKey string) error {
	return connection.Del(SessionKey).Err()
}

//SessionToUserID return the user id stores as value to the session key
func SessionToUserID(SessionKey string) (int, error) {
	key, err := connection.Get(SessionKey).Result()
	if err != nil {
		return 0, err
	}

	if len(key) == 0 {
		return 0, errors.New("No valid user to session key")
	}

	id, err := strconv.Atoi(key)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func refreshSession(Session string) {
	//Value is just the user id
	//The user id is contained inside the sessionkey (Example: session_1_as231fsdf413 -> 1 will be the user id)
	var Value = strings.Split(strings.Split(Session, "session_")[1], "_")[0]
	if err := connection.Set(Session, Value, time.Hour).Err(); err != nil {
		dev.LogError(err, "Couldn't refresh session \""+Session+"\" -> ", err.Error())
	}
}
