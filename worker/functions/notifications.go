package functions

import (
	"encoding/json"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//SendNotifications send the notifications specified in the notification task
func SendNotifications(Task models.WorkerTask) error {
	var Notifications []models.Notification
	err := json.Unmarshal([]byte(Task.Data), &Notifications)
	if err != nil {
		return err
	}

	return nil
}
