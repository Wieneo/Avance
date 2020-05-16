package perms

import "gitlab.gnaucke.dev/tixter/tixter-app/v2/models"

//GetPermissionsToQueue returns the QueuePermissions struct regarding the given user and project
func GetPermissionsToQueue(User models.User, Queue models.Queue) (models.Permissions, *models.QueuePermission, error) {
	perms, err := combinePermissions(User)
	if err != nil {
		return models.Permissions{}, &models.QueuePermission{}, err
	}

	if found, pp := containsQueue(Queue, perms.AccessTo.Queues); found {
		return perms, pp, nil
	}

	return perms, &models.QueuePermission{}, nil
}