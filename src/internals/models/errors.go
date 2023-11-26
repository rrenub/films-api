package models

import "errors"

var ErrNoRecord = errors.New("no mathing records found or empty dataset")
var ErrDuplicatedEntry = errors.New("duplicated entry found")
var ErrInvalidCredentials = errors.New("credentials are invalid")
var ErrTokenExpired = errors.New("access token has expired")
var ErrInvalidToken = errors.New("access token is invalid")
var ErrInvalidAuthHeader = errors.New("Authorization header does not have the correct formatting")
var ErrNotAuthorized = errors.New("User is not authorized to perform this action")
