package cacheable

import (
	"time"
)

func Set(key, value string, expire time.Duration) {
	// if redisMaster.PoolStats().TotalConns > 0 {
	//	if expire == 0 {
	//		redisMaster.Persist(key)
	//	}
	// redisMaster.Set(key, value, expire)
	// tools.HandleError(cmd.Err())
	// } else {
	//	println("There is no redis here")
	//}
}

func Tag(tag, key string) bool {
	// tagKey := new(Tags).key(tag)
	// cmd := redisMaster.LPush(tagKey, key)
	// err := cmd.Err()
	// if err != nil {
	//	return false
	//}
	return true
}

func Expire(key string, expire time.Duration) {
	// redisMaster.Expire(key, expire)
}
