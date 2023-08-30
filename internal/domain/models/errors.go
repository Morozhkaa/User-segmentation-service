// The models package contains a description of the main entities used in the service.
package models

import "fmt"

type ErrorResponse struct {
	ErrorMsg string `json:"error"`
}

var (
	ErrInvalidSlugFormat    = fmt.Errorf("invalid format of parameter 'slug'")    // 400
	ErrInvalidUuidFormat    = fmt.Errorf("invalid format of parameter 'userID'")  // 400
	ErrInvalidPeriodFormat  = fmt.Errorf("invalid format of parameter 'period'")  // 400
	ErrBadRequest           = fmt.Errorf("missing required parameters")           // 400
	ErrSegmentAlreadyExists = fmt.Errorf("segment with this slug already exists") // 400
	ErrSegmentNotFound      = fmt.Errorf("segment not found")                     // 404
)
