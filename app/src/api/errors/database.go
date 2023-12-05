package custom_errors

import (
	"cloudview/app/src/api/middleware/logger"
	"errors"

	"github.com/lib/pq"
)

/*
// Check the gofile to see all ErrorCodes

	type Error struct {
	    Severity         string
	    Code             ErrorCode
	    Message          string
	    Detail           string
	    Hint             string
	    Position         string
	    InternalPosition string
	    InternalQuery    string
	    Where            string
	    Schema           string
	    Table            string
	    Column           string
	    DataTypeName     string
	    Constraint       string
	    File             string
	    Line             string
	    Routine          string
	}
*/
func DBErrors(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		// Handle pq.Error
		if _, found := errorMap[string(pqErr.Code)]; found {
			logger.Logger.Error(pqErr.Message, "code:", pqErr.Code, pqErr.Constraint)
			// Error code found in the map, call the associated function
			return errorMap[string(pqErr.Code)]
		}
	}
	logger.Logger.Error("Unknown Error from DB", err)
	return UnknownError
}

var errorMap = map[string]error{
	"23505": UniqueConstraintViolation,
	"23502": NotNullViolation,
	"23503": ForeignKeyViolation,
	"P0002": NoDataFound,
}

var (
	UniqueConstraintViolation = errors.New("Unique constraint violation")
	NotNullViolation          = errors.New("Not null violation")
	ForeignKeyViolation       = errors.New("Foreign key violation")
	NoDataFound               = errors.New("No data found")
)
