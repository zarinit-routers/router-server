package main

import (
	"flag"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"github.com/zarinit-routers/router-server/pkg/cloud"
	"github.com/zarinit-routers/router-server/pkg/server"
)

func init() {
	flag.BoolFunc("dev-test", "", func(s string) error {
		log.Info("Development/Testing mode enabled")
		viper.Set("dev-test", true)
		return nil
	})
	flag.Parse()
}

func init() {
	viper.SetConfigName("router-config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
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
		viper.SetDefault("server.address", ":8080")
		addr := viper.GetString("server.address")
		if err := r.Run(addr); err != nil {
			log.Fatal("Failed to start HTTP server", "error", err)
		}
	}()

	wg.Wait()
}
