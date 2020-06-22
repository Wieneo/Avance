package worker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.gnaucke.dev/avance/avance-app/v2/worker/functions"

	"gitlab.gnaucke.dev/avance/avance-app/v2/models"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
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
						var Error error
						switch task.Type {
						case models.DeleteUser:
							{
								Error = functions.DeleteUser(task)
								break
							}
						case models.SendNotification:
							{
								retry, err := functions.SendNotifications(task)
								Error = err
								/*We abuse the reoccuring task system here to make tasks execute at a later date
								by setting the lastrun to the current date when scheduling and the interval to the users frequency setting
								*/

								if !retry {
									//Setting the Interval to Invalid here in order to prevent "rescheduling" of the task
									task.Interval.Valid = false
								} else {
									//If an error while connecting to the mailserver happened and the procedure says we can retry -> Set Interval not to false and retry in 300 Seconds
									//This also has the effect that nnew notifications still get appended to the one that couldn't be delivered -> Not flooding the user with e-mails
									dev.LogWarn("A notification delivery failed! Retrying in 300 Seconds...")
									Error = nil
									task.Interval.Int32 = 300
								}
								break
							}
						case models.Debug:
							{
								dev.LogInfo("Debug Task triggered")
								break
							}
						default:
							{
								Error = errors.New("Unknown task type encountered: " + task.Type.String())
								break
							}
						}

						if !task.Interval.Valid {
							if Error != nil {
								dev.LogError(Error, "Task "+strconv.FormatInt(taskid, 10)+" failed: "+Error.Error())
								task.Status = models.Failed
							} else {
								task.Status = models.Finished
							}
						} else {
							if Error != nil {
								dev.LogError(Error, "Reoccuring Task "+strconv.FormatInt(taskid, 10)+" failed: "+Error.Error())
							} else {
								task.Status = models.Idle
								task.LastRun.Valid = true
								task.LastRun.Time = time.Now()
								dev.LogInfo(fmt.Sprintf("Finished Task: (%d) %s -> Next Run in %d Seconds", task.ID, task.Type.String(), task.Interval.Int32))
							}
						}

						if err := db.PatchTask(task); err != nil {
							dev.LogFatal(err, "Couldn't event update task! "+err.Error())
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
