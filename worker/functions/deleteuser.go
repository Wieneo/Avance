package functions

import (
	"encoding/json"
	"fmt"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"

	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//DeleteUser deletes a user from the system
func DeleteUser(Task models.WorkerTask) error {
	var tbd models.User
	err := json.Unmarshal([]byte(Task.Data), &tbd)

	if err != nil {
		return err
	}

	tickets, err := db.GetTicketsOfUser(tbd.ID)
	if err != nil {
		return err
	}

	var Error error

	for _, k := range tickets {
		k.OwnerID.Valid = false
		k.OwnerID.Int64 = 0
		_, err := db.PatchTicket(k)
		if err != nil {
			dev.LogError(err, "Couldn't remove owner from ticket: "+err.Error())
			Error = err
		}
		dev.LogInfo(fmt.Sprintf("[%d] Updated ticket (%d) %s -> Owner set to Nobody", Task.ID, k.ID, k.Title))
	}

	//As the above loop is the last step -> Return the before declared Error
	return Error
}
