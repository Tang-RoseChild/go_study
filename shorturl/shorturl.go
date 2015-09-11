package shorturl

import (
	// "fmt"
	"strconv"
	"sync"
)

var u *URLStore

type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex
}

func (u *URLStore) Get(key string) string {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.urls[key]
}

func (u *URLStore) Set(key, url string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	_, exist := u.urls[key]
	if exist {
		return false
	}

	u.urls[key] = url
	return true
}

func (u *URLStore) Count() int {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return len(u.urls)
}

func (u *URLStore) Put(url string) string {
	key := genKey(u.Count())
	if u.Set(key, url) {
		return key
	}

	return ""
}

func genKey(id int) string {
	return "abc" + strconv.Itoa(id)
}

func NewURLStore() *URLStore {
	if u == nil {
		return &URLStore{urls: make(map[string]string)}
	}

	return u
}
