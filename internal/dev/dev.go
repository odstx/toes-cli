package dev

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Run([]string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go BootRedis(ctx)

	go BootMysql(ctx)

	<-c
	cancel()
}
