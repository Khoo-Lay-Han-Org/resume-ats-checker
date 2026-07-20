package auth_util

import (
	"resuming/database"
	"resuming/database/sqlc"
)

func FindUser(private_id int32) (*sqlc.User, error) {
	return database.FindUser(private_id)
}
