package redis

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"gitlab.gnaucke.dev/avance/avance-app/v2/config"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
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
	_, err := connection.Set("Ping", rand.Intn(100), time.Second).Result()
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
func CreateSession(UserID int64) (string, error) {
	Session := uuid.New()
	if err := connection.Set("session_"+strconv.FormatInt(UserID, 10)+"_"+Session.String(), strconv.FormatInt(UserID, 10), time.Hour).Err(); err != nil {
		return "", err
	}
	return "session_" + strconv.FormatInt(UserID, 10) + "_" + Session.String(), nil
}

//DestroyAllSessions removes all session keys from redis for the specific user
func DestroyAllSessions(UserID int64) error {
	keys, err := connection.Keys("session_" + strconv.FormatInt(UserID, 10) + "_*").Result()
	if err != nil {
		return err
	}

	for _, k := range keys {
		if err := connection.Del(k).Err(); err != nil {
			return err
		}
	}
	return nil
}

//DestroySession removes the session key from redis
func DestroySession(r *http.Request) error {
	session := r.Header.Get("Authorization")
	if len(session) == 0 {
		//Check if maybe cookie was set
		keks, err := r.Cookie("session")
		if err != nil {
			return err
		}

		session = keks.Value
	}

	return connection.Del(session).Err()
}

//SessionToUserID return the user id stores as value to the session key
func SessionToUserID(SessionKey string) (int64, error) {
	key, err := connection.Get(SessionKey).Result()
	if err != nil {
		return 0, err
	}

	if len(key) == 0 {
		return 0, errors.New("No valid user to session key")
	}

	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func refreshSession(Session string) {
	//Value is just the user id
	//The user id is contained inside the sessionkey (Example: session_1_as231fsdf413 -> 1 will be the user id)
	if len(strings.Split(Session, "session_")) == 2 {
		parts := strings.Split(strings.Split(Session, "session_")[1], "_")
		if len(parts) > 1 {
			if err := connection.Expire(Session, time.Hour).Err(); err != nil {
				dev.LogError(err, "Couldn't refresh session \""+Session+"\" -> ", err.Error())
			}
		}
	}

}
