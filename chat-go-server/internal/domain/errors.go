package domain

import "errors"

var ErrNotFound = errors.New("user not found")
var ErrAlreadyExists = errors.New("user with this username already exists")
var ErrChatNotFound = errors.New("chat not found")
var ErrPersonalChatWithoutTargetUser = errors.New("cannot create personal chat without any person")
