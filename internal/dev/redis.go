package dev

import (
	"context"
	"fmt"
	"time"

	rds "github.com/alicebob/miniredis/v2"
)

func BootRedis(ctx context.Context) {
	m := rds.NewMiniRedis()
	err := m.StartAddr(fmt.Sprintf(":%d", 16379))
	if err != nil {
		panic(err)
	}
	defer m.Close()

	for {
		select {
		case <-ctx.Done():
			continue
		default:
			time.Sleep(time.Second * 1)
		}
	}
}
