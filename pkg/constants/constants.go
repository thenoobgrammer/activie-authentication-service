package constants

const (
	MANAGE_ACTIVITIES        = "MANAGE_ACTIVITIES"
	MANAGE_GROUPS            = "MANAGE_GROUPS"
	MANAGE_TEAMS             = "MANAGE_TEAMS"
	REGISTER          string = "REGISTER"
)

const (
	ADMIN = "admin"
	USER  = "user"
)
const (
	HASH_SALT           = 10
	PASSWORD_MIN_LENGTH = 8
)

const (
	ERROR_DURING_TRANSACTION_INIT  string = "error during transaction initialization"
	ERROR_DURING_TRANSACTION       string = "panic on transaction, rolling back"
	ERROR_PANIC_DURING_TRANSACTION string = "error during transaction"
	ERROR_LAST_INSERTED_ID         string = "error retrieving last inserted"
	ERROR_DURING_COMMIT            string = "error during transaction commit"
	ERROR_DURING_QUERY             string = "error during query"
	ERROR_DURING_INSERT            string = "error during insert"
	ERROR_DURING_UPDATE            string = "error during update"
	ERROR_DURING_DELETE            string = "error during deletion"
	ERROR_DURING_ROW_SCAN          string = "error during row scan"
	ERROR_ROWS                     string = "error found in rows scanned"
	WARNING_NO_ROWS_AFFECTED       string = "warning no rows affected"

	ERROR_DURING_CONVERSION string = "error during conversion"

	WARNING_REDIS_USER_CACHING      = "warning there was an issue caching user"
	WARNING_REDIS_USER_INVALIDATION = "warning there was an issue invalidating user"
)

var DEFAULT_PERMISSIONS = [...]string{
	MANAGE_ACTIVITIES,
	MANAGE_GROUPS,
	MANAGE_TEAMS,
	REGISTER,
}
