package db

import (
	"errors"
	"fmt"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//GetTicketRelations returns all relations for the given ticket
func GetTicketRelations(TicketID int64) ([]models.Relation, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting ALL relations for ticket %d", TicketID))
	relations := make([]models.Relation, 0)
	rows, err := Connection.Query(`SELECT "ID", "Ticket1", "Ticket2", "Type" FROM "Relations" WHERE "Ticket1" = $1 OR "Ticket2" = $1`, TicketID)

	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error getting relations for ticket %d -> Returning empty relations array: %s", TicketID, err.Error()))
		return make([]models.Relation, 0), err
	}

	for rows.Next() {
		var singleRelation models.Relation

		var ticket1, ticket2 int64
		var relationType models.RelationType
		rows.Scan(&singleRelation.ID, &ticket1, &ticket2, &relationType)

		//GetTicket is used with false ResolveRelations property to prevent possibly resolving hundrets of ticket relations

		if ticket1 == TicketID {
			dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Resolving related ticket %d", TicketID, ticket2))
			ticket, _, err := GetTicketUnsafe(ticket2, models.WantedProperties{})

			if err != nil {
				dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Error getting related ticket %d -> Returning empty relations array: %s", TicketID, ticket2, err.Error()))
				return make([]models.Relation, 0), err
			}

			singleRelation.OtherTicket = ticket
			singleRelation.Type = relationType
		} else {
			dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Resolving related ticket %d", TicketID, ticket1))
			ticket, _, err := GetTicketUnsafe(ticket1, models.WantedProperties{})

			if err != nil {
				dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Error getting related ticket %d -> Returning empty relations array: %s", TicketID, ticket1, err.Error()))
				return make([]models.Relation, 0), err
			}

			singleRelation.OtherTicket = ticket

			//Invert meanings of relation
			switch relationType {
			case models.References:
				{
					singleRelation.Type = models.ReferencedBy
					break
				}
			case models.ReferencedBy:
				{
					singleRelation.Type = models.References
					break
				}
			case models.ParentOf:
				{
					singleRelation.Type = models.ChildOf
					break
				}
			case models.ChildOf:
				{
					singleRelation.Type = models.ParentOf
					break
				}
			default:
				{
					err := errors.New("Unknown Relation Type")
					dev.LogError(err, "Unknown relation type: "+string(relationType))
					return make([]models.Relation, 0), err
				}
			}
		}

		relations = append(relations, singleRelation)
	}

	dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Got %d relations", TicketID, len(relations)))

	return relations, nil
}

//AddRelation adds a relation to a ticket
func AddRelation(Ticket1, Ticket2 int64, Type models.RelationType) (int64, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Adding relation from ticket %d to %d with type %d", Ticket1, Ticket2, Type))
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Relations" ("Ticket1", "Ticket2", "Type") VALUES ($1, $2, $3) RETURNING "ID"`, Ticket1, Ticket2, Type).Scan(&newID)
	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Refreshing Ticket1 (%d) of relation %d", Ticket1, newID))
		err = TicketWasModified(Ticket1)
		if err == nil {
			dev.LogDebug(fmt.Sprintf("[DB] Refreshing Ticket2 (%d) of relation %d", Ticket2, newID))
			err = TicketWasModified(Ticket2)
		}
	}

	dev.LogDebug(fmt.Sprintf("[DB] Created relation %d between %d <-> %d (%d)", newID, Ticket1, Ticket2, Type))
	return newID, err
}

//DeleteRelation deletes a relation to a ticket
func DeleteRelation(ID int64) error {
	dev.LogDebug(fmt.Sprintf("[DB] Deleting relation %d", ID))
	var ticket1, ticket2 int64
	err := Connection.QueryRow(`DELETE FROM "Relations" WHERE "ID" = $1 RETURNING "Ticket1", "Ticket2"`, ID).Scan(&ticket1, &ticket2)
	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Refreshing Ticket1 (%d) of relation %d", ticket1, ID))
		err = TicketWasModified(ticket1)
		if err == nil {
			dev.LogDebug(fmt.Sprintf("[DB] Refreshing Ticket2 (%d) of relation %d", ticket2, ID))
			err = TicketWasModified(ticket2)
		}
	}
	return err
}
