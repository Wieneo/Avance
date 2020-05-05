package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
)

//SampleConfig templates the used config
type SampleConfig struct {
	Debug    bool
	Port     int
	Postgres struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
	}
	Redis struct {
		Host     string
		Port     int
		Database int
		Sentinel struct {
			Enabled   bool
			Master    string
			Endpoints []string
		}
	}
}

//CurrentConfig stores the currently used config
var CurrentConfig SampleConfig

//LoadConfig tries to load the current config from ENV / File
func LoadConfig() {
	EnvEnabled := false
	var Error error
	//Try to get ENV Variables
	if len(os.Getenv("TIX_ENV")) != 0 {
		EnvEnabled, Error = strconv.ParseBool(os.Getenv("TIX_ENV"))
		if Error != nil {
			EnvEnabled = false
			dev.LogInfo("TIX_ENV seems not to be a boolean! Please check your variables.")
		}
	}

	if EnvEnabled {
		dev.LogInfo("Using ENV-Variables for configuration")

		CurrentConfig.Debug, Error = strconv.ParseBool(os.Getenv("TIX_Debug"))
		if Error != nil {
			dev.LogWarn("Debug is not a boolean! Using false as default")
			CurrentConfig.Debug = false
		}
		dev.DeubgLogging = CurrentConfig.Debug

		//Parse listen port to int
		CurrentConfig.Port, Error = strconv.Atoi(os.Getenv("TIX_Port"))
		if Error != nil {
			dev.LogWarn("Port is not a number! Using 8000 as default")
			CurrentConfig.Port = 8000
		}

		CurrentConfig.Postgres.Host = os.Getenv("TIX_Postgres_Host")
		CurrentConfig.Postgres.Username = os.Getenv("TIX_Postgres_Username")
		CurrentConfig.Postgres.Password = os.Getenv("TIX_Postgres_Password")
		CurrentConfig.Postgres.Database = os.Getenv("TIX_Postgres_Database")

		//Parse Postgres Port to int
		CurrentConfig.Postgres.Port, Error = strconv.Atoi(os.Getenv("TIX_Postgres_Port"))
		if Error != nil {
			dev.LogWarn("Postgres port is not a number! Using 5432 as default")
			CurrentConfig.Postgres.Port = 5432
		}

		CurrentConfig.Redis.Host = os.Getenv("TIX_Redis_Host")
		//Parse Redis Database to int
		CurrentConfig.Redis.Port, Error = strconv.Atoi(os.Getenv("TIX_Redis_Database"))
		if Error != nil {
			dev.LogWarn("Redis database is not a number! Using 0 as default")
			CurrentConfig.Redis.Database = 0
		}

		//Parse Redis Port to int
		CurrentConfig.Redis.Port, Error = strconv.Atoi(os.Getenv("TIX_Redis_Port"))
		if Error != nil {
			dev.LogWarn("Redis port is not a number! Using 6379 as default")
			CurrentConfig.Redis.Port = 6379
		}

		CurrentConfig.Redis.Sentinel.Enabled, Error = strconv.ParseBool(os.Getenv("TIX_Redis_Sentinel_Enabled"))
		if Error != nil {
			dev.LogWarn("Redis Sentinel is not a boolean! Using false as default")
			CurrentConfig.Redis.Sentinel.Enabled = false
		}

		if CurrentConfig.Redis.Sentinel.Enabled {
			CurrentConfig.Redis.Sentinel.Master = os.Getenv("TIX_Redis_Sentinel_Master")

			//Search for sentinels
			for _, e := range os.Environ() {
				if strings.HasPrefix(e, "TIX_Redis_Sentinel_Endpoint_") {
					value := strings.Split(e, "=")[1]
					CurrentConfig.Redis.Sentinel.Endpoints = append(CurrentConfig.Redis.Sentinel.Endpoints, value)
				}
			}

		}

	} else {
		dir, _ := os.Getwd()
		dev.LogInfo("Using ", dir+"/config/config.json", " for configuration")
		rawBytes, err := ioutil.ReadFile(dir + "/config/config.json")
		if err != nil {
			dev.LogFatal("Couldn't read config:", err.Error())
			return
		}
		err = json.Unmarshal(rawBytes, &CurrentConfig)
		if err != nil {
			dev.LogFatal("Couldn't read config:", err.Error())
			return
		}
		dev.DeubgLogging = CurrentConfig.Debug

	}

	if CurrentConfig.Redis.Sentinel.Enabled {
		dev.LogInfo(len(CurrentConfig.Redis.Sentinel.Endpoints), "Redis Sentinel Nodes configured")
		dev.LogInfo("Looking for", CurrentConfig.Redis.Sentinel.Master, "as master")
	}

	dev.LogInfo("Config was read completely!")
}
