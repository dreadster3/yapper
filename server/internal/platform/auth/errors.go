package auth

import "errors"

var ErrForbidden = errors.New("forbidden to access this resource")
