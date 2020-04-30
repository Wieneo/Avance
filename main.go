package main

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/redis"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/server"

	//Imported to be used with database/sql
	_ "github.com/lib/pq"
)

const (
	//MainVersion is set to the current version
	MainVersion = "0.1"
	//Channel is set to the current dev channel
	Channel = "Dev"
)

func main() {
	dev.LogInfo("Starting App Version", MainVersion, Channel)
	config.LoadConfig()
	db.Init()
	redis.Init()
	server.HTTPInit()
}
