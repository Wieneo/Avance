package models

import "fmt"

//User is the default user struct from the database
type User struct {
	ID          int64
	Username    string
	Mail        string
	Permissions Permissions
	Firstname   string
	Lastname    string
	Settings    UserSettings
}

//UserSettings struct stores all settings about a single user
type UserSettings struct {
	EnabledNotificationChannels []NotificationChannel
	//EnabledNotificationChannelsReadable represents the human readable equivalent of EnabledNotificationChannels. You cannot update the settings via this struct
	EnabledNotificationChannelsReadable []string
	//NotificationFrequency defines the frequency in which notifications should be send (in seconds)
	NotificationFrequency       int
	NotificationAboutNewTickets bool
	NotificationAboutUpdates    bool
	NotificationAfterInvolvment bool
}

//NotificationChannel defines the type of a notification
type NotificationChannel int

const (
	//Mail is currently the only supported one
	Mail NotificationChannel = iota
)

func (e NotificationChannel) String() string {
	switch e {
	case Mail:
		return "Mail"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}
