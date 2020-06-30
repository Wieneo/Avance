package perms

import (
	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
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

		if k.Permissions.CanCreateUsers {
			PermissionSet.CanCreateUsers = true
		}

		if k.Permissions.CanCreateGroups {
			PermissionSet.CanCreateGroups = true
		}

		if k.Permissions.CanModifyUsers {
			PermissionSet.CanModifyUsers = true
		}

		if k.Permissions.CanModifyGroups {
			PermissionSet.CanModifyGroups = true
		}

		if k.Permissions.CanDeleteUsers {
			PermissionSet.CanDeleteUsers = true
		}

		if k.Permissions.CanDeleteGroups {
			PermissionSet.CanDeleteGroups = true
		}

		if k.Permissions.CanChangePermissionsGlobal {
			PermissionSet.CanChangePermissionsGlobal = true
		}

		if k.Permissions.CanSeeWorker {
			PermissionSet.CanSeeWorker = true
		}

		if k.Permissions.CanChangeWorker {
			PermissionSet.CanChangeWorker = true
		}

		for _, k := range k.Permissions.AccessTo.Projects {
			if found, project := containsProjectPermission(k, PermissionSet.AccessTo.Projects); !found {
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
				if k.CanCreateSeverities {
					project.CanCreateSeverities = true
				}
				if k.CanModifySeverities {
					project.CanModifySeverities = true
				}
				if k.CanRemoveSeverities {
					project.CanRemoveSeverities = true
				}
				if k.CanCreateStatuses {
					project.CanCreateStatuses = true
				}
				if k.CanModifyStatuses {
					project.CanModifyStatuses = true
				}
				if k.CanRemoveStatuses {
					project.CanRemoveStatuses = true
				}
			}
		}

		for _, k := range k.Permissions.AccessTo.Queues {
			if found, project := containsQueuePermission(k, PermissionSet.AccessTo.Queues); !found {
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

func containsProject(project models.Project, projects []models.ProjectPermission) (bool, *models.ProjectPermission) {
	for i, k := range projects {
		if k.ProjectID == project.ID {
			return true, &projects[i]
		}
	}

	return false, nil
}

func containsProjectPermission(project models.ProjectPermission, projects []models.ProjectPermission) (bool, *models.ProjectPermission) {
	for i, k := range projects {
		if k.ProjectID == project.ProjectID {
			return true, &projects[i]
		}
	}

	return false, nil
}

func containsQueuePermission(queue models.QueuePermission, queues []models.QueuePermission) (bool, *models.QueuePermission) {
	for i, k := range queues {
		if k.QueueID == queue.QueueID {
			return true, &queues[i]
		}
	}

	return false, nil
}

func permsContainQueue(queue models.Queue, queues []models.QueuePermission) (bool, *models.QueuePermission) {
	for i, k := range queues {
		if k.QueueID == queue.ID {
			return true, &queues[i]
		}
	}

	return false, nil
}
