package main

import (
	"flag"
	"sync"

	"github.com/spf13/viper"
	"github.com/zarinit-routers/router-server/pkg/cloud"
	"github.com/zarinit-routers/router-server/pkg/server"
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

	// start cloud connection loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		cloud.ServeConnection()
	}()

	// start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		r := server.New()
		addr := viper.GetString("server.address")
		if addr == "" {
			addr = ":8080"
		}
		_ = r.Run(addr)
	}()

	wg.Wait()
}
