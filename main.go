package main

import (
	"os"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/redis"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/server"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/worker"

	//Imported to be used with database/sql
	_ "github.com/lib/pq"
)

const usageMessage = "Please specify if the container should start as App (--app) or as Worker (--worker)"

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "--app":
			{
				dev.LogInfo("Starting App Version", utils.AppVersion, utils.AppChannel)
				config.LoadConfig()
				db.Init()
				redis.Init()
				server.HTTPInit()
				break
			}
		case "--worker":
			{
				dev.LogInfo("Starting Worker Version", utils.WorkerVersion, utils.WorkerChannel)
				config.LoadConfig()
				db.Init()
				worker.StartQueueService()
				worker.StartServing()
				break
			}
		default:
			{
				dev.LogInfo(usageMessage)
				break
			}
		}
	} else {
		dev.LogInfo(usageMessage)
	}

}
