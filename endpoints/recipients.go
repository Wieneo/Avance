package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/templates"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
)

//AddRecipient adds a recipient to the ticket
func AddRecipient(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)
	ticketid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[8], 10, 64)

	var req []recipientRequest

	rawBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(rawBytes, &req)
	if err != nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "JSON body is malformed: "+err.Error())
		return
	}

	queue, found, err := db.GetQueue(projectid, queueid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.QueueNoPerms)
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	allperms, perms, err := perms.GetPermissionsToQueue(user, queue)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !allperms.Admin && !perms.CanEditTicket {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to edit tickets in that queue.")
		return
	}

	ticket, found, err := db.GetTicket(ticketid, queueid, models.WantedProperties{Recipients: true})
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.TicketNotFound)
		return
	}

	var created []int64

	//Error-Check first
	for i, k := range req {
		if utils.IsEmpty(k.Mail) && utils.IsEmpty(k.User) {
			w.WriteHeader(406)
			dev.ReportUserError(w, "Error at "+strconv.Itoa(i)+": Either User or Mail must be specified!")
			return
		}
		if !utils.IsEmpty(k.Mail) && !utils.IsEmpty(k.User) {
			w.WriteHeader(406)
			dev.ReportUserError(w, "Error at "+strconv.Itoa(i)+": Only a user OR a mail address can be specified!")
			return
		}

		if _, err := strconv.Atoi(k.Type.String()); err == nil {
			w.WriteHeader(406)
			dev.ReportUserError(w, "Error at "+strconv.Itoa(i)+": Unknown type specified!")
			return
		}

		if !utils.IsEmpty(k.User) {
			userid, found, err := db.SearchUser(k.User)
			if err != nil {

				utils.ReportInternalErrorToUser(err, w)
				return
			}

			if !found {
				w.WriteHeader(404)
				dev.ReportUserError(w, "Error at "+strconv.Itoa(i)+": Specified user wasn't found")
				return
			}

			for _, exist := range recipientsWithType(ticket.Recipients, k.Type) {
				if exist.User.Valid {
					if exist.User.Value.ID == userid {
						w.WriteHeader(406)
						dev.ReportUserError(w, "Error at "+strconv.Itoa(i)+": Specified user already is a recipient of that type!")
						return
					}
				}
			}

			if k.Type == models.Readers {
				found := false
				for _, exist := range recipientsWithType(ticket.Recipients, models.Admins) {
					if exist.User.Valid {
						if exist.User.Value.ID == userid {
							found = true
						}
					}
				}

				if !found {
					for _, exist := range recipientsWithType(ticket.Recipients, models.Requestors) {
						if exist.User.Valid {
							if exist.User.Value.ID == userid {
								found = true
							}
						}
					}
				}

				if found {
					w.WriteHeader(406)
					dev.ReportUserError(w, "Error at "+strconv.Itoa(i)+": Specified user already is a requestor/admin")
					return
				}
			} else {
				for _, exist := range recipientsWithType(ticket.Recipients, models.Readers) {
					if exist.User.Valid {
						if exist.User.Value.ID == userid {
							w.WriteHeader(406)
							dev.ReportUserError(w, "Error at "+strconv.Itoa(i)+": Specified user already is a reader")
							return
						}
					}
				}
			}

			//Dummy recipient
			newRecipient := models.Recipient{
				User: struct {
					Valid bool
					Value models.User
				}{
					Valid: true,
					Value: models.User{
						ID: userid,
					},
				},
			}
			switch k.Type {
			case models.Admins:
				{
					ticket.Recipients.Admins = append(ticket.Recipients.Admins, newRecipient)
					break
				}
			case models.Requestors:
				{
					ticket.Recipients.Requestors = append(ticket.Recipients.Requestors, newRecipient)
					break
				}
			case models.Readers:
				{
					ticket.Recipients.Readers = append(ticket.Recipients.Readers, newRecipient)
					break
				}
			}

		} else {
			for _, exist := range recipientsWithType(ticket.Recipients, k.Type) {
				if !exist.User.Valid {
					if strings.ToLower(exist.Mail) == strings.ToLower(k.Mail) {
						w.WriteHeader(406)
						dev.ReportUserError(w, "Error at "+strconv.Itoa(i)+": Specified e-mail already is a recipient of that type!")
						return
					}
				}
			}

			if k.Type == models.Readers {
				found := false
				for _, exist := range recipientsWithType(ticket.Recipients, models.Admins) {
					if !exist.User.Valid {
						if strings.ToLower(exist.Mail) == strings.ToLower(k.Mail) {
							found = true
						}
					}
				}

				if !found {
					for _, exist := range recipientsWithType(ticket.Recipients, models.Requestors) {
						if !exist.User.Valid {
							if strings.ToLower(exist.Mail) == strings.ToLower(k.Mail) {
								found = true
							}
						}
					}
				}

				if found {
					w.WriteHeader(406)
					dev.ReportUserError(w, "Error at "+strconv.Itoa(i)+": Specified e-mail already is a requestor/admin")
					return
				}
			} else {
				for _, exist := range recipientsWithType(ticket.Recipients, models.Readers) {
					if !exist.User.Valid {
						if strings.ToLower(exist.Mail) == strings.ToLower(k.Mail) {
							w.WriteHeader(406)
							dev.ReportUserError(w, "Error at "+strconv.Itoa(i)+": Specified e-mail already is a reader")
							return
						}
					}
				}
			}

			//Dummy recipient
			newRecipient := models.Recipient{
				Mail: k.Mail,
			}
			switch k.Type {
			case models.Admins:
				{
					ticket.Recipients.Admins = append(ticket.Recipients.Admins, newRecipient)
					break
				}
			case models.Requestors:
				{
					ticket.Recipients.Requestors = append(ticket.Recipients.Requestors, newRecipient)
					break
				}
			case models.Readers:
				{
					ticket.Recipients.Readers = append(ticket.Recipients.Readers, newRecipient)
					break
				}
			}
		}
	}
	//ToDo: Rollback if something fails here!
	//Execute changes
	for _, k := range req {
		if !utils.IsEmpty(k.User) {
			userid, _, err := db.SearchUser(k.User)
			if err != nil {

				utils.ReportInternalErrorToUser(err, w)
				return
			}

			newid, err := db.AddUserRecipient(ticket.ID, userid, k.Type)
			if err != nil {

				utils.ReportInternalErrorToUser(err, w)
				return
			}

			created = append(created, newid)
		} else {
			newid, err := db.AddMailRecipient(ticket.ID, k.Mail, k.Type)
			if err != nil {

				utils.ReportInternalErrorToUser(err, w)
				return
			}

			created = append(created, newid)
		}
	}

	json.NewEncoder(w).Encode(struct {
		Recipients []int64
	}{
		created,
	})

}

func recipientsWithType(Recipients models.RecipientCollection, Type models.RecipientType) []models.Recipient {
	switch Type {
	case models.Requestors:
		return Recipients.Requestors
	case models.Admins:
		return Recipients.Admins
	case models.Readers:
		return Recipients.Readers
	}

	return make([]models.Recipient, 0)
}

func allRecipients(Ticket models.Ticket) []models.Recipient {
	all := append(Ticket.Recipients.Requestors, Ticket.Recipients.Admins...)
	all = append(all, Ticket.Recipients.Readers...)
	return all
}

//DeleteRecipient removes a recipient from the ticket
func DeleteRecipient(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)
	ticketid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[8], 10, 64)
	recipientid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[10], 10, 64)

	queue, found, err := db.GetQueue(projectid, queueid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.QueueNotFound)
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	allperms, perms, err := perms.GetPermissionsToQueue(user, queue)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !allperms.Admin && !perms.CanEditTicket {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to edit tickets in that queue.")
		return
	}

	ticket, found, err := db.GetTicket(ticketid, queueid, models.WantedProperties{Recipients: true})
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.TicketNotFound)
		return
	}

	found = false
	for _, k := range allRecipients(ticket) {
		if k.ID == recipientid {
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.RecipientNotFound)
		return
	}

	if err := db.RemoveRecipient(recipientid); err != nil {

		utils.ReportInternalErrorToUser(err, w)
	} else {
		json.NewEncoder(w).Encode(struct {
			Recipient string
		}{
			fmt.Sprintf("Recipient %d deleted", recipientid),
		})
	}
}
