package mongo

import (
	"regexp"

	"github.com/iqdf/benjerry-service/domain"
)

var (
	// RegexpNoDocumentError ...
	RegexpNoDocumentError = regexp.MustCompile(`^mongo: no documents in result$`)

	// RegexpDuplicateError ...
	RegexpDuplicateError = regexp.MustCompile(`duplicate key error collection: .+ index:(?P<Field>.+) dup key`)
)

func checkNoDocumentError(dbError error) bool {
	return RegexpNoDocumentError.Match([]byte(dbError.Error()))
}

func checkWriteDuplicateError(dbError error) bool {
	return RegexpDuplicateError.Match([]byte(dbError.Error()))
}

// TranslateError converts mongo DB error into
// approriate application errors
func TranslateError(dbError error) error {
	if dbError == nil {
		return nil
	}

	switch {
	case checkNoDocumentError(dbError):
		return domain.ErrResourceNotFound

	case checkWriteDuplicateError(dbError):
		return domain.ErrConflict
	}

	return domain.ErrInternalServerError
}
