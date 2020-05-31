package worker

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/worker/functions"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
)

//StartQueueService starts the main thread for queue execution
func StartQueueService() {
	//Don't block main thread
	go func() {
		for {
			active, err := db.GetWorkerStatus()
			if err != nil {
				dev.LogError(err, "Couldn't check for worker status: "+err.Error())
				time.Sleep(time.Second * 5)
				continue
			}

			if active {
				taskid, err := db.ReserveTask()
				if err != nil {
					if strings.Contains(err.Error(), "converting NULL to int64 is unsupported") {
						dev.LogInfo("No tasks available. Sleeping for 5 seconds")
						time.Sleep(time.Second * 5)
					} else {
						dev.LogError(err, err.Error())
					}
				} else {
					task, err := db.GetTask(taskid)
					if err != nil {
						dev.LogError(err, "Couldn't pickup task "+strconv.FormatInt(taskid, 10)+": "+err.Error())
						task.Status = models.Failed
						if db.PatchTask(task) != nil {
							dev.LogFatal(err, "Couldn't event update task! "+err.Error())
						}
					} else {
						dev.LogInfo(fmt.Sprintf("Picked Up Task: (%d) %s %s", task.ID, task.Type.String(), task.QueuedAt.String()))
						//Everything is fine
						switch task.Type {
						case models.DeleteUser:
							{
								err := functions.DeleteUser(task)
								if err != nil {
									dev.LogError(err, "Task "+strconv.FormatInt(taskid, 10)+" failed: "+err.Error())
									task.Status = models.Failed
								} else {
									task.Status = models.Finished
								}

								if db.PatchTask(task) != nil {
									dev.LogFatal(err, "Couldn't event update task! "+err.Error())
								}
								break
							}
						}
					}
				}
			} else {
				dev.LogInfo("Worker is not active! Sleeping 30 Seconds.")
				time.Sleep(time.Second * 30)
			}
		}
	}()
}
