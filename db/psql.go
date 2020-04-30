package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
)

var connection *sql.DB

//Init is called after config is read to connect to postgres
func Init() {
	dev.LogInfo("Postgres is being initialized")
	//FROM: https://godoc.org/github.com/lib/pq
	//ToDo: Postgres Port is currently ignored!
	connStr := fmt.Sprint("postgres://", config.CurrentConfig.Postgres.Username, ":", config.CurrentConfig.Postgres.Password, "@", config.CurrentConfig.Postgres.Host, ":", config.CurrentConfig.Postgres.Port, "/", config.CurrentConfig.Postgres.Database, "?sslmode=disable")
	dev.LogDebug("Connecting to ", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		dev.LogFatal("Couldn't initialize Postgres: ", err)
		//Program will terminate here
	}
	connection = db
	dev.LogInfo("-- Connection to Postgres established --")
	migrate()
}

func migrate() {
	dev.LogInfo("Preparing database to be migrated...")

	//Get current verison from "Version" table
	var currentVersion int
	rows, err := connection.Query(`SELECT "Schema" FROM "Version"`)
	if err != nil {
		dev.LogFatal("Prepare failed:", err)
	}
	if !rows.Next() {
		dev.LogFatal("Version Table is empty! Something went horribly wrong... Check your database.")
	}

	//Scan returns error
	if rows.Scan(&currentVersion) != nil {
		dev.LogFatal("Version Table is malformed! Something went horribly wrong... Check your database.")
	}

	cwd, _ := os.Getwd()
	files, err := ioutil.ReadDir(fmt.Sprint(cwd, "/db/migrations/"))
	if err != nil {
		dev.LogFatal("Couldn't find/read migrations: ", err.Error())
	}

	var newestVersion int
	migrationsAvailable := make(map[int]string, 0)
	for _, k := range files {
		version, err := strconv.Atoi(strings.Split(k.Name(), "-")[0])
		if err != nil {
			dev.LogDebug("Ignoring non-migration file:", k.Name())
		}

		if version > newestVersion {
			newestVersion = version
		}
		migrationsAvailable[version] = k.Name()
	}

	dev.LogInfo("Database:", currentVersion, " | Code:", newestVersion)
	if currentVersion < newestVersion {
		for currentVersion < newestVersion {
			currentVersion++
			dev.LogInfo("Applying", migrationsAvailable[currentVersion])
			rawBytes, err := ioutil.ReadFile(cwd + "/db/migrations/" + migrationsAvailable[currentVersion])
			if err != nil {
				dev.LogFatal("Couldn't read migration:", err.Error())
			}

			_, err = connection.Exec(string(rawBytes))
			if err != nil {
				dev.LogFatal("Couldn't apply migration:", err.Error())
			}

			_, err = connection.Exec(`UPDATE "Version" SET "Schema" = $1`, currentVersion)
			if err != nil {
				dev.LogFatal("Couldn't apply migration:", err.Error())
			}
		}
	} else {
		dev.LogInfo("No migrations needed!")
	}
}
