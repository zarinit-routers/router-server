package main

import (
	"flag"
	"sync"

	"github.com/spf13/viper"
	"github.com/zarinit-routers/router-server/pkg/cloud"
)

func init() {
	flag.BoolFunc("dev-test", "", func(s string) error {
		viper.Set("dev-test", true)
		return nil
	})
	flag.Parse()
}

func init() {
	viper.SetConfigName("router-config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
}

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		cloud.ServeConnection()
	}()

	wg.Wait()
}
