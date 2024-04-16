package errs

import "errors"

var ErrUserExists = errors.New("user already exists")
var ErrInvalidLoginOrPassword = errors.New("invalid login or password")
var ErrInvalidLoginOrAuthToken = errors.New("invalid login or auth token")