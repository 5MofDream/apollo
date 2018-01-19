package cache

import (
	"time"
	"sync"
	"fmt"
	"encoding/gob"
	"os"
	"io"
)
const (
	// 没有过期时间标志
	NoExpiration time.Duration = -1
	//默认的过期时间
	DefaultExpiration time.Duration = 0
)

type Cache struct {
	defaultExpiration time.Duration
	items             map[string]Item //缓存数据项存储在map中
	mu                sync.RWMutex    //读写锁
	gcInterval        time.Duration   //过期数据项清理周期
	stopGc            chan bool
}

// 过期缓存数据项清理
func (c *Cache) gcLoop() {
	ticker := time.NewTicker(c.gcInterval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.stopGc:
			ticker.Stop()
			return
		}
	}
}

func (c *Cache) del(k string) {
	delete(c.items, k)
}

func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			c.del(k)
		}
	}
}

func (c *Cache) Set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[k] = Item{
		Object:     v,
		Expiration: e,
	}
}

// 有序关系的集合
func (c *Cache) ZSet(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.items[k]
	if !found {
		c.items[k] = Item{
			Object:     v,
			Expiration: e,
		}
	}


}


func (c *Cache) set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = Item{
		Object:     v,
		Expiration: e,
	}
}

func (c *Cache) get(k string) (interface{}, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	if item.Expired() {
		return nil, false
	}
	return item.Object, true
}

func (c *Cache) Add(k string, v interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s aleady exists", k)
	}
	c.set(k, v, d)
	c.mu.Unlock()
	return nil
}

func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expired() {
		return nil, false
	}
	c.mu.RUnlock()
	return item.Object, true
}

func (c *Cache) Replace(k string, v interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	c.set(k, v, d)
	c.mu.Unlock()
	return nil
}

func (c *Cache) Delete(k string) {
	c.mu.Lock()
	c.del(k)
	c.mu.Unlock()
}

func (c *Cache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, v := range c.items {
		gob.Register(v.Object)
	}
	err = enc.Encode(&c.items)
	return
}

func (c *Cache) SaveToFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	if err = c.Save(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func (c *Cache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := map[string]Item{}
	err := dec.Decode(&items)
	if err == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		for k, v := range items {
			ov, found := c.items[k]
			if !found || ov.Expired() {
				c.items[k] = v
			}
		}
	}
	return err
}

func (c *Cache) LoadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	if err = c.Load(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func (c *Cache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = map[string]Item{}
}

func (c *Cache) StopGc() {
	c.stopGc <- true
}

// 初始化一个新的cache服务 公开出去
func NewCache() *Cache {
	//i:= config.Read("cache.defaultExpiration") TODO 走配置
	c := &Cache{
		defaultExpiration: 1,
		gcInterval:        1,
		items:             map[string]Item{},
		stopGc:            make(chan bool),
	}
	go c.gcLoop()
	return c
}
