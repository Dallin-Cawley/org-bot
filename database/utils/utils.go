package utils

import (
	"strings"
)

const (
	NO_ROWS_RETURNED = "no rows returned from database"
)

// NoRowsRead determines if the database error is reporting that no record was returned.
func NoRowsRead(err error) bool {
	return strings.Contains(err.Error(), NO_ROWS_RETURNED)
}
