package cacheable

import (
	"github.com/go-redis/redis/v8"
)

type cache struct {
	client redis.Cmdable
}

func SetUpRedis(db int, Cluster int) {

}

func SetUpRedisCluster(nodes []string) error {
	// redisMaster = redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs: nodes,
	// })
	// redisSlave = redisMaster
	// if redisMaster == nil || redisSlave == nil {
	//	return fmt.Errorf("redisCluster: nil")
	//}
	return nil
}
