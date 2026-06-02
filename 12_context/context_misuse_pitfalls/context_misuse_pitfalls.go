package main

import (
	"context"
	"fmt"
	"time"
)

type BadService struct {
	ctx context.Context
}

func newBadService(ctx context.Context) *BadService {
	return &BadService{ctx: ctx}
}

func goodOperation(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}

func main() {
	parent, cancel := context.WithCancel(context.Background())
	service := newBadService(parent)

	cancel()

	fmt.Println("bad pattern: storing context in struct means stale/canceled context can persist")
	fmt.Println("stored context err:", service.ctx.Err())

	freshCtx, freshCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer freshCancel()

	err := goodOperation(freshCtx)
	fmt.Println("good operation err:", err)
}
