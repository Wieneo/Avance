package models

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
	Notification NotificationSettings
}

//NotificationSettings is nested in UserSettings and stores all information about sending notifications to the user
type NotificationSettings struct {
	MailNotificationEnabled bool
	//NotificationFrequency defines the frequency in which notifications should be send (in seconds)
	NotificationFrequency       int
	NotificationAboutNewTickets bool
	NotificationAboutUpdates    bool
	NotificationAfterInvolvment bool
}
