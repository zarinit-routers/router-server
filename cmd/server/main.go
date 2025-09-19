package main

import (
	"flag"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"github.com/zarinit-routers/router-server/internal/user"
	"github.com/zarinit-routers/router-server/pkg/cloud"
	"github.com/zarinit-routers/router-server/pkg/server"
	"github.com/zarinit-routers/router-server/pkg/storage"
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
	viper.AddConfigPath("/etc/zarinit/")
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error while reading config file", "error", err)
	}
	viper.AutomaticEnv()

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {

	if err := storage.Check(); err != nil {
		log.Fatal("Key-value storage is not available, check failed", "error", err)
	}

	if err := user.EnsureCreated(); err != nil {
		log.Fatal("Error while checking default user", "error", err)
	}

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
		log.Info("Starting HTTP server", "address", addr)
		if err := r.Run(addr); err != nil {
			log.Fatal("Failed to start HTTP server", "error", err)
		}
	}()

	wg.Wait()
}
