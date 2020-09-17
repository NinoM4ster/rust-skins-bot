package main

import (
	"context"
	"time"
)

func newCtx(t time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), t*time.Second)
}
