package cacheable

import "context"

func Del(ctx context.Context, key string) {
	// redisMaster.Del(ctx, key)
}

func DelByTag(ctx context.Context, tag string) {
	// tagsKey := new(Tags).key(tag)
	// var tags Tags
	// tags = GetKeys(ctx, tag)
	// for _, index := range tags {
	//	redisMaster.Del(ctx, index)
	//}
	// redisMaster.Del(ctx, tagsKey)
}
