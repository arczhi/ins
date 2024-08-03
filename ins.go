package ins

import (
	"errors"
	"hash/crc32"
	"math"
	"sync"
	"time"
)

func New() *Ins {
	return &Ins{
		buckets: [16]bucket{},
	}
}

type Ins struct {
	buckets [16]bucket
}

type bucket struct {
	db map[string]item
	sync.Mutex
}

type item struct {
	value     interface{}
	expiredAt int64
}

// set key value
func (i *Ins) Set(key string, value interface{}) error {
	p := i.partition(key)
	i.buckets[p].Lock()
	defer i.buckets[p].Unlock()
	i.set(p, key, item{value: value, expiredAt: -1})
	return nil
}

func (i *Ins) SetNx(key string, val interface{}) error {
	p := i.partition(key)
	i.buckets[p].Lock()
	defer i.buckets[p].Unlock()
	dbVal, ok := i.get(p, key)
	if ok && dbVal.expiredAt != -1 && time.Now().After(time.Unix(dbVal.expiredAt, 0)) {
		i.del(p, key)
	}
	if ok {
		return errors.New("key exists")
	}
	i.set(p, key, item{value: val, expiredAt: -1})
	return nil
}

func (i *Ins) SetNxEx(key string, val interface{}, exp int64) error {
	p := i.partition(key)
	i.buckets[p].Lock()
	defer i.buckets[p].Unlock()
	_, ok := i.get(p, key)
	if ok {
		return errors.New("key exists")
	}

	i.set(p, key, item{value: val, expiredAt: time.Now().Unix() + exp})
	return nil
}

func (i *Ins) SetEx(key string, val interface{}, exp int64) error {
	p := i.partition(key)
	i.buckets[p].Lock()
	defer i.buckets[p].Unlock()
	i.set(p, key, item{value: val, expiredAt: time.Now().Unix() + exp})
	return nil
}

func (i *Ins) set(partition int, key string, val item) {
	if i.buckets[partition].db == nil {
		i.buckets[partition].db = make(map[string]item)
	}
	i.buckets[partition].db[key] = val
}

// return nil,false if key not found
func (i *Ins) Get(key string) (interface{}, bool) {
	p := i.partition(key)
	val, ok := i.get(p, key)
	if ok {
		if val.expiredAt != -1 && time.Now().After(time.Unix(val.expiredAt, 0)) {
			i.del(p, key)
			return nil, false
		}
	}
	return val.value, ok
}

func (i *Ins) get(partition int, key string) (item, bool) {
	if i.buckets[partition].db == nil {
		return item{}, false
	}
	val, ok := i.buckets[partition].db[key]
	return val, ok
}

func (i *Ins) Del(key string) error {
	p := i.partition(key)
	return i.del(p, key)
}

func (i *Ins) del(partition int, key string) error {
	i.buckets[partition].Lock()
	defer i.buckets[partition].Unlock()
	delete(i.buckets[partition].db, key)
	return nil
}

func (i *Ins) partition(key string) int {
	return int(math.Mod(float64(crc32.ChecksumIEEE([]byte(key))), 16))
}
