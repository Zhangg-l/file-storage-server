package redis

import (
	"context"
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	cli := RedisPool()
	if cli == nil {
		cli = newRedisPool()
	}
	ctx := context.Background()
	status := cli.Set(ctx, "k1", "zhangg", 0)
	fmt.Println(status.Args()...)
	res := cli.Get(ctx, "k1")
	fmt.Println(res.Result())
}
