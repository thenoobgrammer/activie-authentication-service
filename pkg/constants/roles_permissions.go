package constants

const (
	ACTIVITIES_VIEW        string = "activities-view"
	ACTIVITIES_CREATE      string = "activities-create"
	ACTIVITIES_DELETE      string = "activities-delete"
	ACTIVITIES_REGISTER    string = "activities-register"
	GROUP_CREATE           string = "group-create"
	GROUP_DELETE           string = "group-delete"
	GROUP_SEARCH           string = "group-search"
	USERS_VIEW_ALL         string = "see-all-users"
	USERS_VIEW_DETAIL      string = "see-detail-users"
	SEND_MESSAGES          string = "send-messages"
	NOTIFICATIONS_RECEIVE  string = "notifications-receive"
	GOOGLE_LOCATION_SEARCH string = "google-location-search"
	MAP_VIEW               string = "map-view"

	ADMIN = "admin"
	USER  = "user"
)

var DEFAULT_PERMISSIONS = [...]string{
	ACTIVITIES_VIEW,
	ACTIVITIES_CREATE,
	ACTIVITIES_DELETE,
	NOTIFICATIONS_RECEIVE,
	GOOGLE_LOCATION_SEARCH,
	MAP_VIEW,
}
