package cacheable

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
)

const KeyPrefix = "Cache"

var EmptyCacheError = errors.New("Cachable: No results found")

type Tags []string

type Cache struct {
	key    string
	value  string
	expire time.Duration
}
type ICacheable interface {
	GetKey() string
	Get() error
	Set(expire time.Duration) error
}
type Cacheable struct {
	LastUpdated time.Time `json:"last_updated,omitempty"`
	prefix      string
	cacheID     string
	parent      ICacheable
}

func (tag Tags) key(key string) string {
	return "redis_tags_" + key
}
func NewCachable(prefix, id string, parentPtr ICacheable) (*Cacheable, error) {
	if reflect.ValueOf(parentPtr).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("parent field in cachable must be a pointer")
	}
	return &Cacheable{prefix: prefix, cacheID: id, parent: parentPtr}, nil
}

func (c *Cacheable) GetKey() string {
	return KeyPrefix + c.prefix + ":" + c.cacheID
}

func (c *Cacheable) Get(ctx context.Context) error {
	r := Get(ctx, c.parent.GetKey())
	if r != "" {
		err := json.Unmarshal([]byte(r), c.parent)
		if err != nil {
			return err
		}
		return nil
	}
	return EmptyCacheError
}

func (c *Cacheable) Set(expire time.Duration) error {
	c.LastUpdated = time.Now()
	v, err := json.Marshal(c.parent)
	if err != nil {
		return err
	}
	Set(c.parent.GetKey(), string(v), expire)
	return nil
}
