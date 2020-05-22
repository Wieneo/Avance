package db

import (
	"errors"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetTicketRelations returns all relations for the given ticket
func GetTicketRelations(TicketID int64) ([]models.Relation, error) {
	relations := make([]models.Relation, 0)
	rows, err := Connection.Query(`SELECT "ID", "Ticket1", "Ticket2", "Type" FROM "Relations" WHERE "Ticket1" = $1 OR "Ticket2" = $1`, TicketID)

	if err != nil {
		return make([]models.Relation, 0), err
	}

	for rows.Next() {
		var singleRelation models.Relation

		var ticket1, ticket2 int64
		var relationType models.RelationType
		rows.Scan(&singleRelation.ID, &ticket1, &ticket2, &relationType)

		//GetTicket is used with false ResolveRelations property to prevent possibly resolving hundrets of ticket relations

		if ticket1 == TicketID {
			ticket, _, err := GetTicket(ticket2, false)

			if err != nil {
				return make([]models.Relation, 0), err
			}

			singleRelation.OtherTicket = ticket
			singleRelation.Type = relationType
		} else {
			ticket, _, err := GetTicket(ticket1, false)

			if err != nil {
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

	return relations, nil
}

//AddRelation adds a relation to a ticket
func AddRelation(Ticket1, Ticket2 int64, Type models.RelationType) (int64, error) {
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Relations" ("Ticket1", "Ticket2", "Type") VALUES ($1, $2, $3) RETURNING "ID"`, Ticket1, Ticket2, Type).Scan(&newID)
	return newID, err
}

//DeleteRelation deletes a relation to a ticket
func DeleteRelation(ID int64) error {
	_, err := Connection.Exec(`DELETE FROM "Relations" WHERE "ID" = $1`, ID)
	return err
}
