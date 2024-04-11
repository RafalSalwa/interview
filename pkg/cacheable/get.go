package cacheable

import "context"

func Get(ctx context.Context, key string) string {
	// if redisSlave.PoolStats().TotalConns > 0 {

	// val, err := redisSlave.Get(ctx, key).Result()
	// if err != nil {
	//	if err.Error() != "redis: nil" {
	//		return GetFromFallBack(ctx, key)
	//	}
	//}
	// return val
	// } else {
	//	return GetFromFallBack(key)
	//}
	return ""
}

func GetFromFallBack(ctx context.Context, key string) string {
	// var val string
	// if redisFallback.PoolStats().TotalConns > 0 {
	// if redisFallback == nil {
	//	return ""
	//}
	// val, err := redisFallback.Get(ctx, key).Result()
	// if err != nil {
	//	if err.Error() != "redis: nil" {
	//		return GetFromFallBack(ctx, key)
	//	}
	//	return val
	//}
	//
	//}
	// return val
	return ""
}

func GetList(ctx context.Context, key string, start, stop int64) ([]string, error) {
	// return redisSlave.LRange(ctx, key, start, stop).Result()
	return []string{""}, nil
}

func GetKeys(ctx context.Context, tag string) []string {
	var tags Tags
	var err error
	tagsKey := tags.key(tag)
	tags, err = GetList(ctx, tagsKey, 0, 10)
	if err != nil {
		return nil
	}
	return tags
}
