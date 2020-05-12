package main

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/redis"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/server"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"

	//Imported to be used with database/sql
	_ "github.com/lib/pq"
)

func main() {
	dev.LogInfo("Starting App Version", utils.MainVersion, utils.Channel)
	config.LoadConfig()
	db.Init()
	redis.Init()
	server.HTTPInit()
}
