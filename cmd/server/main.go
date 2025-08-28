package main

import (
	"sync"

	"github.com/zarinit-routers/router-server/pkg/cloud"
)

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		cloud.ServeConnection()
	}()
}
