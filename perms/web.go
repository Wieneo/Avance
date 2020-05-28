package perms

import "gitlab.gnaucke.dev/tixter/tixter-app/v2/models"

//IsAdmin returns true if the user gets admin privileges from a group / itself
func IsAdmin(User models.User) (bool, error) {
	perms, err := CombinePermissions(User)
	if err != nil {
		return false, err
	}

	return perms.Admin, nil
}
