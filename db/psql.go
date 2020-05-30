package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//Connection stores the current connection to postgres
var Connection *sql.DB

//Migration stores a single database migration
type Migration struct {
	Name    string
	Targets []string
	After   string
}

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

	rows, err := Connection.Query(`SELECT "Name" FROM "Patches"`)
	if err != nil {
		dev.LogError(err, err.Error())
		//If table doesn't exist
		if strings.Contains(err.Error(), "does not exist") {
			//Fix for old instances
			if _, err := Connection.Query(`SELECT "Schema" FROM "Version"`); err == nil {
				dev.LogInfo("Old Database System detected -> Migrating")

				//Create needed patches table
				if _, err := Connection.Exec(`CREATE TABLE public."Patches" ("Name" text NOT NULL);`); err != nil {
					dev.LogFatal(err, "Couldn't migrate to new migration system! "+err.Error())
				} else {
					migrate()
					return
				}
			}

			deploy()
			migrate()

			userid, _ := CreateUser(models.User{
				Username:  "Admin",
				Firstname: "The",
				Lastname:  "Admin",
				Mail:      "root@localhost",
				Permissions: models.Permissions{
					Admin: true,
				},
			}, "tixter")

			//The following is used to make debugging and developing the APP easier when used with Gitlab Auto DevOPS
			//Detect if deployed via GITLAB
			if len(os.Getenv("GITLAB_ENVIRONMENT_NAME")) > 0 {
				dev.LogInfo("Instance was deployed via Gitlab. Deploying example data")
				projectid, _ := CreateProject("Auto DevOPS", "Default project created by Gitlab Auto DevOPS")
				qid, _ := CreateQueue("Development", projectid)
				statusid, _ := CreateStatus(true, "Open", "green", true, projectid)
				severityid, _ := CreateSeverity(true, "Normal", "green", 10, projectid)
				ticket1, _ := CreateTicket("Pipeline broken", "My pipeline is broken!", qid, true, 0, severityid, statusid, false, "")
				ticket2, _ := CreateTicket("Create User", "Please create a user for my new staff member!", qid, false, userid, severityid, statusid, false, "")
				AddRelation(ticket1, ticket2, models.ParentOf)
				AddRelation(ticket2, ticket1, models.ReferencedBy)
				AddAction(ticket1, models.Comment, "Comment was added", "This is the first comment on this new instance!", userid)
				AddAction(ticket1, models.Answer, "Answer was added", "This is the first answer on this new instance!", userid)
				AddAction(ticket2, models.Answer, "Answer was added", "This is the second Answer on this new instance!<br>Even with <i>formatting</i>", userid)

			}
			return
		}

		dev.LogFatal(err, err.Error())
	}
	/*LEGACY!

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
			//Ignore versions that arent mapped
			if len(migrationsAvailable[currentVersion]) > 0 {
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
		}
	} else {
		dev.LogInfo("No migrations needed!")
	}
	*/
	//appliedPatches contains all patches from the "Patches" database table. These are tracked via their names
	appliedPatches := make([]string, 0)
	//localPatches contains all patches from the ./db/migrations/ folder. All files with the suffix .migrate.json will be included
	localPatches := make([]Migration, 0)

	//SELECT resides on the second line of this function. Is also used to check if the table even exists
	for rows.Next() {
		var patch string
		rows.Scan(&patch)
		appliedPatches = append(appliedPatches, patch)
	}

	rows.Close()

	cwd, _ := os.Getwd()
	files, err := ioutil.ReadDir(fmt.Sprint(cwd, "/db/migrations/"))
	if err != nil {
		dev.LogFatal(err, "Couldn't find/read migrations: ", err.Error())
	}

	//Read all json files
	for _, k := range files {
		if strings.HasSuffix(k.Name(), ".migrate.json") {
			var migration Migration
			rawBytes, err := ioutil.ReadFile(fmt.Sprint(cwd, "/db/migrations/", k.Name()))
			if err != nil {
				dev.LogFatal(err, "Couldn't read migration file: "+err.Error())
			}

			if err := json.Unmarshal(rawBytes, &migration); err != nil {
				dev.LogFatal(err, "Error in migration "+k.Name()+": "+err.Error())
			}

			//Name, Targets and After shouldn't be empty
			if len(migration.Name) == 0 || len(migration.Targets) == 0 || len(migration.After) == 0 {
				dev.LogWarn("Skipping empty migration: " + migration.Name + "(" + k.Name() + ")")
				continue
			}

			//Check if the migration or a migration with that name already was included
			found := false
			for _, patch := range localPatches {
				if patch.Name == migration.Name {
					found = true
				}
			}

			//App should exit if migration with duplicate name exists.
			if found {
				dev.LogFatal(errors.New("Duplicate Migration"), "Duplicate Migartion-Name: "+migration.Name)
			}

			localPatches = append(localPatches, migration)
		}
	}

	dev.LogInfo(fmt.Sprintf("Database has %d Patches applied. Local Patches available: %d", len(appliedPatches), len(localPatches)))

	//Check what patches already were applied to the database
	neededPatches := make([]Migration, 0)
	for _, k := range localPatches {
		found := false
		for _, db := range appliedPatches {
			if k.Name == db {
				found = true
			}
		}

		if !found {
			neededPatches = append(neededPatches, k)
		}
	}

	dev.LogInfo(fmt.Sprintf("%d Migration(s) need to be applied", len(neededPatches)))

	//MAIN APPLY
	for len(neededPatches) > 0 {
		patchesProcessed := 0
		for i, k := range neededPatches {

			//Check if dependency of migartion already was applied
			baseExists := false
			//{{BASE}} is a placeholder. Basically means the patch has no dependency and can installed right after the base image was imported
			if k.After != "{{BASE}}" {
				for _, db := range appliedPatches {
					if db == k.After {
						baseExists = true
					}
				}
			} else {
				baseExists = true
			}

			//If dependency was applied
			if baseExists {
				dev.LogInfo("Applying " + k.Name)
				sqlstring := "BEGIN;"

				//Targets contains all SQL Files wich are included in this migration
				//Concat these here
				for _, k := range k.Targets {
					rawBytes, err := ioutil.ReadFile(cwd + "/db/migrations/" + k)

					if err != nil {
						dev.LogFatal(err, "Couldn't read target:", err.Error())
					}

					//Little fix to maybe catch some BEGIN; END; statements which break the transaction
					//If BEGIN; / END; is specified in the SQL file, a single target file can fail and others be still applied.
					//Thats not good because next start the whole migration will be applied again (even with the already applied target)
					rawString := strings.ReplaceAll(strings.ReplaceAll(string(rawBytes), "BEGIN;", ""), "END;", "")

					sqlstring += rawString
				}

				sqlstring += "END;"
				if _, err = Connection.Exec(sqlstring); err != nil {
					dev.LogFatal(err, "Couldn't apply migration:", err.Error())
				}

				//Add Migration to applied patches (database and local)
				if _, err := Connection.Exec(`INSERT INTO "Patches" VALUES ($1)`, k.Name); err != nil {
					dev.LogFatal(err, "Couldn't apply migration:", err.Error())
				}

				appliedPatches = append(appliedPatches, k.Name)

				//Removed this migration from neededPAtches
				if len(neededPatches) >= i+2 {
					neededPatches = append(neededPatches[:i], neededPatches[i+1:]...)
				} else {
					//No preceeding item left
					neededPatches = make([]Migration, 0)
				}

				//Track how many migrations where applied in this run (multiple runs can occur because of dependencies)
				patchesProcessed++
			}
		}

		if patchesProcessed == 0 {
			errorMessage := "No patches could be applied! Maybe you have a dependency loop / missing dependency?\nThe following patches remain to be applied: "
			for _, k := range neededPatches {
				rawBytes, _ := json.Marshal(k)
				errorMessage += "\n" + string(rawBytes)
			}

			dev.LogFatal(errors.New("Patches couldn't be applied"), errorMessage)
		}
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

	if err != nil {
		dev.LogFatal(err, "Couldn't create administrator!", err.Error())
	}

}
