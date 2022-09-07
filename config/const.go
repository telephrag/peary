package config

const (
	// Used for giving name to database and log files.
	// If changed won't be consistent with container and executable name.
	NAME = "peary"

	// Used for giving name to part of database that stores data about players who
	// took role.
	DB_PLAYERS_BUCKET_NAME = "players"
)
