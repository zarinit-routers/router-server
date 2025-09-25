package config

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func TestGetConfiguration(t *testing.T) {
	viper.Set("wifi-hotspot.interface", "wlan0")
	config, err := FromFile("../../.././test/data/hostapd.conf")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config)

	if err := config.SetNecessaryDefaultOptions(); err != nil {
		t.Fatalf("failed apply necessary defaults: %s", err)
	}

	if err := config.SetPassphrase("some password with spaces"); err != nil {
		log.Fatalf("failed set passphrase: %s", err)
	}

	if err := config.Write(); err != nil {
		t.Fatal(err)
	}

}
