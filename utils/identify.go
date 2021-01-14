package utils

import "sync"

type uuid struct {
	mu   sync.Mutex
	uuid uint64
}

func (u *uuid) Get() uint64 {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.uuid = u.uuid + 1
	return u.uuid
}

var Uuid uuid
