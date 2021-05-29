package db

import "sync"

var StateToUser = make(map[string]string)
var StateMapLock = &sync.RWMutex{}
var UserToToken = make(map[string]string)
var UserMapLock = &sync.RWMutex{}
