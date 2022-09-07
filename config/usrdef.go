package config

// For advanced users
const (
	// Time interval in seconds. Each time interval passes, database is scanned
	// for players whose role is expired and deletes their records.
	DB_CHANGESTREAM_SLEEP_SECONDS = 60

	// Maximum time interval for handling command from player to the bot.
	// Note that this must not exceed docker's wait time on `docker stop %container%`
	// (default is 10 seconds) since some command handlers might not complete in time
	// before container with the bot stops.
	CMD_HANDLER_TIMEOUT_SECONDS = 5
)

const (
	// Name of the role players will take.
	BOT_ROLE_NAME = "Want to play"

	// Color of the role the players take.
	// Take hex RGB value of color you want and convert it to integer.
	// For example, default value bellow is 4AF47 in hex.
	// Make it stand out among other roles on your server.
	BOT_ROLE_COLOR = 307015
)
