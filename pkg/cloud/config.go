package cloud

import (
	"fmt"
	"net/url"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

var (
	config *viper.Viper
)

func init() {
	config = viper.New()
	config.SetConfigName("cloud-config")
	config.AddConfigPath(".")
	config.AddConfigPath("/etc/")

	// Defaults
	config.SetDefault("reconnect-timeout", time.Second*10)

	// Read

	if err := config.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("There is no config file", "err", err)
		} else {
			log.Error("There was an error reading the config file", "err", err)
		}

	}
}

type ConnectionConfig struct {
	Passphrase     string
	OrganizationId string
	CloudHost      string
	WSPort         int
}

func (c *ConnectionConfig) GetWebsocketURL() string {
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", c.CloudHost, c.WSPort),
		Path:   CloudWebSocketPath,
	}
	return u.String()
}
func (c *ConnectionConfig) GetAuthURL() string {
	u := url.URL{
		Scheme: "http",
		Host:   c.CloudHost,
		Path:   CloudAuthPath,
	}
	return u.String()
}

func (c *ConnectionConfig) validate() error {
	if c.Passphrase == "" {
		return fmt.Errorf("passphrase is empty")
	}
	if c.OrganizationId == "" {
		return fmt.Errorf("organization id is empty")
	}
	if c.CloudHost == "" {
		return fmt.Errorf("cloud host is not set")
	}
	if c.WSPort == 0 {
		return fmt.Errorf("websocket port is not set")
	}
	if c.WSPort < 0 {
		return fmt.Errorf("invalid websocket port %d", c.WSPort)
	}
	return nil
}

func getCloudConfig() (*ConnectionConfig, error) {
	c := ConnectionConfig{
		Passphrase:     config.GetString("passphrase"),
		OrganizationId: config.GetString("organization-id"),
		CloudHost:      config.GetString("cloud-host"),
		WSPort:         config.GetInt("websocket-port"),
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	log.Info("Cloud config", "config", c)

	return &c, nil
}
