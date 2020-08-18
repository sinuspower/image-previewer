package cache //nolint:golint,stylecheck

import (
	"crypto/sha1" //nolint:go-lint
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear() error
	Cap() int
	GetFile(string) ([]byte, bool, error)
	PutFile(string, []byte) error
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	path     string // path to cache dir in filesystem
	queue    List
	items    map[Key]*listItem
	mutex    *sync.Mutex
}

func NewCache(capacity int, path string) (Cache, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0700); err != nil {
			return nil, err
		}
	}

	return &lruCache{
		capacity: capacity,
		path:     path,
		queue:    NewList(),
		items:    make(map[Key]*listItem),
		mutex:    &sync.Mutex{},
	}, nil
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mutex.Lock()
	if itm, ok := lc.items[key]; ok { // refresh
		refreshed := cacheItem{key, value}
		itm.Value = refreshed
		lc.queue.MoveToFront(lc.items[key])
		lc.mutex.Unlock()

		return true
	}
	// insert
	lc.items[key] = lc.queue.PushFront(cacheItem{key, value})
	if len(lc.items) > lc.capacity { // remove old record
		old := lc.queue.Back()
		key := old.Value.(cacheItem).key
		delete(lc.items, key)
		lc.queue.Remove(old)
		// delete file if exists
		fileName := lc.path + "/" + string(key)
		if _, err := os.Stat(fileName); err == nil {
			_ = os.Remove(fileName)
		}
	}
	lc.mutex.Unlock()

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mutex.Lock()
	if itm, ok := lc.items[key]; ok { // return value
		lc.queue.MoveToFront(lc.items[key])
		lc.mutex.Unlock()

		return itm.Value.(cacheItem).value, true
	}
	lc.mutex.Unlock()

	return nil, false
}

func (lc *lruCache) Clear() error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	lc.queue = NewList()
	lc.items = make(map[Key]*listItem)
	if err := os.RemoveAll(lc.path); err != nil {
		return err
	}

	return nil
}

func (lc *lruCache) Cap() int {
	return lc.capacity
}

func (lc *lruCache) GetFile(path string) ([]byte, bool, error) {
	key := getHash(path)
	if _, ok := lc.Get(Key(key)); ok { // hashed filename is in cache
		f, err := os.Open(lc.path + "/" + key)
		if err != nil {
			return nil, false, err
		}
		defer f.Close()
		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, false, err
		}

		return bytes, true, nil
	}

	return nil, false, nil
}

func (lc *lruCache) PutFile(path string, data []byte) error {
	key := getHash(path)
	fileName := lc.path + "/" + key

	err := ioutil.WriteFile(fileName, data, 0600)
	if err != nil {
		return err
	}

	if lc.Set(Key(key), 0) {
		return errors.New("already in cache, rewritten")
	}

	return nil
}

func getHash(source string) string {
	sha1Bytes := sha1.Sum([]byte(source)) //nolint:go-lint

	return hex.EncodeToString(sha1Bytes[:])
}
