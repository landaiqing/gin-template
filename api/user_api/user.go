package user_api

import "sync"

type UserAPI struct{}

var mu sync.Mutex
