package helper

import (
	"go_auth-project/dto"
	"sync"
	"time"
)

type TokenData struct {
	AccessToken  string
	RefreshToken string
	User         dto.UserResponse
	ExpiresAt    time.Time
}

var (
	store = make(map[string]TokenData)
	mu    sync.Mutex
)

func SetMemory(code string, data TokenData, ttl time.Duration) {
	mu.Lock()
	defer mu.Unlock()

	data.ExpiresAt = time.Now().Add(ttl)
	store[code] = data
}

func Getmemory(code string) (TokenData, bool) {
	mu.Lock()
	defer mu.Unlock()

	data, ok := store[code]
	if !ok || time.Now().After(data.ExpiresAt) {
		delete(store, code)
		return TokenData{}, false
	}

	delete(store, code) // one-time use
	return data, true
}
