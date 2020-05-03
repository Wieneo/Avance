package perms

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//CombinePermissions combines the permission from the user and the assigned groups to produce the highest permission set
func CombinePermissions(User models.User) (models.Permissions, error) {
	var PermissionSet models.Permissions
	PermissionSet = User.Permissions

	groups, err := db.GetGroups(User)
	if err != nil {
		return PermissionSet, err
	}

	for _, k := range groups {
		//If on group makes user an admin
		if k.Permissions.Admin {
			PermissionSet.Admin = true
			break
		}

		for _, k := range k.Permissions.AccessTo.Projects {
			if found, project := containsProject(k, PermissionSet.AccessTo.Projects); !found {
				PermissionSet.AccessTo.Projects = append(PermissionSet.AccessTo.Projects, k)
			} else {
				if k.CanSee {
					project.CanSee = true
				}
				if k.CanModify {
					project.CanModify = true
				}
				if k.CanChangePermissions {
					project.CanChangePermissions = true
				}
				if k.CanCreateQueues {
					project.CanCreateQueues = true
				}
				if k.CanModifyQueues {
					project.CanModifyQueues = true
				}
				if k.CanRemoveQueues {
					project.CanRemoveQueues = true
				}
			}
		}

		for _, k := range k.Permissions.AccessTo.Queues {
			if found, project := containsQueue(k, PermissionSet.AccessTo.Queues); !found {
				PermissionSet.AccessTo.Queues = append(PermissionSet.AccessTo.Queues, k)
			} else {
				if k.CanSee {
					project.CanSee = true
				}
				if k.CanModify {
					project.CanModify = true
				}
				if k.CanChangePermissions {
					project.CanChangePermissions = true
				}
				if k.CanCreateTicket {
					project.CanCreateTicket = true
				}
				if k.CanEditTicket {
					project.CanEditTicket = true
				}
			}
		}
	}

	return PermissionSet, nil
}

func containsProject(project models.ProjectPermission, projects []models.ProjectPermission) (bool, *models.ProjectPermission) {
	for i, k := range projects {
		if k.ProjectID == project.ProjectID {
			return true, &projects[i]
		}
	}

	return false, nil
}

func containsQueue(queue models.QueuePermission, queues []models.QueuePermission) (bool, *models.QueuePermission) {
	for i, k := range queues {
		if k.QueueID == queue.QueueID {
			return true, &queues[i]
		}
	}

	return false, nil
}
