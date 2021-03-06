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

	"gitlab.gnaucke.dev/avance/avance-app/v2/config"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//Connection stores the current connection to postgres
var Connection *sql.DB

//Migration stores a single database migration
type Migration struct {
	Name    string
	Targets []string
	After   string
}

const migrationsPath = "/db/migrations/"

//Init is called after config is read to connect to postgres
func Init(ApplyMigrations bool) {
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

	migrate(ApplyMigrations)
}

func migrate(ApplyMigrations bool) {
	dev.LogInfo("Preparing database to be migrated...")

	rows, err := Connection.Query(`SELECT "Name" FROM "Patches"`)
	if err != nil {
		//dev.LogError(err, err.Error())
		//If table doesn't exist
		if strings.Contains(err.Error(), "does not exist") && ApplyMigrations {
			//Fix for old instances
			if _, err := Connection.Query(`SELECT "Schema" FROM "Version"`); err == nil {
				dev.LogInfo("Old Database System detected -> Migrating")

				//Create needed patches table
				if _, err := Connection.Exec(`CREATE TABLE public."Patches" ("Name" text NOT NULL);`); err != nil {
					dev.LogFatal(err, "Couldn't migrate to new migration system! "+err.Error())
				} else {
					migrate(ApplyMigrations)
					return
				}
			}

			deploy()
			migrate(ApplyMigrations)

			userid, _ := CreateUser(models.User{
				Username:  "Admin",
				Firstname: "The",
				Lastname:  "Admin",
				Mail:      "root@localhost",
				Permissions: models.Permissions{
					Admin: true,
				},
				Settings: models.UserSettings{
					Notification: models.NotificationSettings{
						MailNotificationAboutNewTickets: true,
						MailNotificationAboutUpdates:    true,
						MailNotificationAfterInvolvment: false,
						MailNotificationEnabled:         true,
						MailNotificationFrequency:       300,
					},
				},
			}, "avance")

			//The following is used to make debugging and developing the APP easier when used with Gitlab Auto DevOPS
			//Detect if deployed via GITLAB
			if len(os.Getenv("GITLAB_ENVIRONMENT_NAME")) > 0 || config.CurrentConfig.SetupDemo {
				if !config.CurrentConfig.SetupDemo {
					dev.LogInfo("Instance was deployed via Gitlab. Deploying example data")
				} else {
					dev.LogInfo("SetupDemo flag has been set! Deploying example data")
				}
				projectid, _ := CreateProject("Auto DevOPS", "Default project created by Gitlab Auto DevOPS")
				qid, _ := CreateQueue("Development", projectid)
				statusid, _ := CreateStatus(true, "Open", "green", true, projectid)
				severityid, _ := CreateSeverity(true, "Normal", "green", 10, projectid)

				severity, _, _ := GetSeverityUNSAFE(severityid)
				status, _, _ := GetStatusUNSAFE(statusid)

				newTicket1 := models.CreateTicket{
					Title:         "Pipeline broken",
					Description:   "My pipeline is broken!",
					Queue:         qid,
					OwnedByNobody: true,
					Severity:      severity,
					Status:        status,
					IsStalled:     false,
				}

				newTicket2 := models.CreateTicket{
					Title:         "Create User",
					Description:   "Please create a user for my new staff member!",
					Queue:         qid,
					OwnedByNobody: true,
					Severity:      severity,
					Status:        status,
					IsStalled:     false,
				}

				ticket1, _ := CreateTicket(newTicket1)
				ticket2, _ := CreateTicket(newTicket2)
				AddRelation(ticket1, ticket2, models.ParentOf)
				AddRelation(ticket2, ticket1, models.ReferencedBy)

				user, _, _ := GetUser(userid)

				AddAction(ticket1, models.Comment, "Comment was added", "This is the first comment on this new instance!", models.Issuer{Valid: true, Issuer: user})
				AddAction(ticket1, models.Answer, "Answer was added", "This is the first answer on this new instance!", models.Issuer{Valid: true, Issuer: user})
				AddAction(ticket2, models.Answer, "Answer was added", "This is the second Answer on this new instance!<br>Even with <i>formatting</i>", models.Issuer{Valid: true, Issuer: user})
			}
			return
		}

		if !ApplyMigrations {
			dev.LogInfo("Database Schema isn't deployed. Waiting for it...")
			stallMigration(ApplyMigrations)
			return
		}

		dev.LogFatal(err, err.Error())
	}

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
	files, err := ioutil.ReadDir(fmt.Sprint(cwd, migrationsPath))
	if err != nil {
		dev.LogFatal(err, "Couldn't find/read migrations: ", err.Error())
	}

	//Read all json files
	for _, k := range files {
		if strings.HasSuffix(k.Name(), ".migrate.json") {
			var migration Migration
			rawBytes, err := ioutil.ReadFile(fmt.Sprint(cwd, migrationsPath, k.Name()))
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

	if !ApplyMigrations && len(neededPatches) > 0 {
		dev.LogInfo("There are still unapplied patches. Waiting for App to apply them!")
		stallMigration(ApplyMigrations)
		return
	}

	//MAIN APPLY
	for len(neededPatches) > 0 {
		patchesProcessed := 0

		patchesForThisRun := make([]Migration, 0)
		for _, k := range neededPatches {
			patchesForThisRun = append(patchesForThisRun, k)
		}
		for _, k := range patchesForThisRun {
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
					rawBytes, err := ioutil.ReadFile(cwd + migrationsPath + k)

					if err != nil {
						dev.LogFatal(err, "Couldn't read target:", err.Error())
					}

					//Little fix to maybe catch some BEGIN; END; statements which break the transaction
					//If BEGIN; / END; is specified in the SQL file, a single target file can fail and others be still applied.
					//Thats not good because next start the whole migration will be applied again (even with the already applied target)
					//rawString := strings.ReplaceAll(strings.ReplaceAll(string(rawBytes), "BEGIN;", ""), "END;", "")

					sqlstring += string(rawBytes)
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

				//Removed this migration from neededPatches
				for mi, mk := range neededPatches {
					if mk.Name == k.Name {
						if len(neededPatches) >= mi+1 {
							neededPatches = append(neededPatches[:mi], neededPatches[mi+1:]...)
							break
						} else {
							neededPatches = neededPatches[:mi]
							break
						}
					}
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
	rawBytes, err := ioutil.ReadFile(cwd + migrationsPath + "base.sql")
	if err != nil {
		dev.LogFatal(err, "Couldn't read "+cwd+migrationsPath+"base.sql! Please check permissions and roles!")
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

func stallMigration(ApplyMigrations bool) {
	time.Sleep(5 * time.Second)
	migrate(ApplyMigrations)
}
