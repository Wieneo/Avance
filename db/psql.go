package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//Connection stores the current connection to postgres
var Connection *sql.DB

//Init is called after config is read to connect to postgres
func Init() {
	dev.LogInfo("Postgres is being initialized")
	//FROM: https://godoc.org/github.com/lib/pq
	connStr := "ERROR"
	if len(config.CurrentConfig.Postgres.ConnectionString) == 0 {
		connStr = fmt.Sprint("postgres://", config.CurrentConfig.Postgres.Username, ":", config.CurrentConfig.Postgres.Password, "@", config.CurrentConfig.Postgres.Host, ":", config.CurrentConfig.Postgres.Port, "/", config.CurrentConfig.Postgres.Database, "?sslmode=disable")
	} else {
		connStr = config.CurrentConfig.Postgres.ConnectionString + "?sslmode=disable"
		dev.LogInfo("DATABASE_URL was set! Using that instead of TIX_Postgres_*")
	}

	dev.LogDebug("Connecting to ", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		dev.LogFatal(err, "Couldn't initialize Postgres: ", err)
		//Program will terminate here
	}
	Connection = db
	migrate()
}

func migrate() {
	dev.LogInfo("Preparing database to be migrated...")

	//Get current verison from "Version" table
	var currentVersion int
	rows, err := Connection.Query(`SELECT "Schema" FROM "Version"`)
	if err != nil {
		dev.LogError(err, err.Error())
		//If table doesn't exist
		if strings.Contains(err.Error(), "does not exist") {
			deploy()
			migrate()
			return
		}

		dev.LogFatal(err, err.Error())
	}
	if !rows.Next() {
		dev.LogFatal(err, "Version Table is empty! Something went horribly wrong... Check your database.")
	}

	//Scan returns error
	if rows.Scan(&currentVersion) != nil {
		dev.LogFatal(err, "Version Table is malformed! Something went horribly wrong... Check your database.")
	}

	cwd, _ := os.Getwd()
	files, err := ioutil.ReadDir(fmt.Sprint(cwd, "/db/migrations/"))
	if err != nil {
		dev.LogFatal(err, "Couldn't find/read migrations: ", err.Error())
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
				dev.LogFatal(err, "Couldn't read migration:", err.Error())
			}

			_, err = Connection.Exec(string(rawBytes))
			if err != nil {
				dev.LogFatal(err, "Couldn't apply migration:", err.Error())
			}

			_, err = Connection.Exec(`UPDATE "Version" SET "Schema" = $1`, currentVersion)
			if err != nil {
				dev.LogFatal(err, "Couldn't apply migration:", err.Error())
			}
		}
	} else {
		dev.LogInfo("No migrations needed!")
	}
}

func deploy() {
	dev.LogInfo("Trying automatic deploy...")
	dev.LogInfo("Looking for other/remaining tables")
	rows, err := Connection.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name;")
	if err != nil {
		dev.LogFatal(err, "Couldn't read schema! Please check permissions and roles!")
	}

	if rows.Next() {
		dev.LogFatal(err, "Schema not empty! Not deploying due to safety concerns")
	}

	rows.Close()

	cwd, _ := os.Getwd()
	rawBytes, err := ioutil.ReadFile(cwd + "/db/migrations/base.sql")
	if err != nil {
		dev.LogFatal(err, "Couldn't read "+cwd+"/db/migrations/base.sql! Please check permissions and roles!")
	}

	_, err = Connection.Query(string(rawBytes))
	if err != nil {
		dev.LogFatal(err, "Couldn't deploy schema!", err.Error())
	}

	dev.LogInfo("Schema deployed. Waiting 5 Seconds until restarting migration process")
	time.Sleep(5 * time.Second)

	dev.LogInfo("Setting default schema version")
	_, err = Connection.Exec(`INSERT INTO "Version" VALUES ('0')`)
	if err != nil {
		dev.LogFatal(err, "Couldn't set schema version:", err.Error())
	}

	_, err = CreateUser(models.User{
		Username:  "Admin",
		Firstname: "Admin",
		Lastname:  "istrator",
		Mail:      "root@localhost",
		Permissions: models.Permissions{
			Admin: true,
		},
	}, "tixter")

	if err != nil {
		dev.LogFatal(err, "Couldn't create administrator!", err.Error())
	}

}
