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
	config.AddConfigPath("/etc/zarinit/")

	// Defaults
	config.SetDefault("reconnect-timeout", time.Second*10)

	// Read

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("There is no config file", "err", err)
		} else {
			log.Error("There was an error reading the config file", "err", err)
		}

	}
}

const (
	CloudWebSocketPath = "/api/ipc/connect"
	CloudAuthPath      = "/api/organizations/authorize-node"
)

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

type ValidationError struct {
	err error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", e.err.Error())
}

func NewValidationError(err error) error {
	return &ValidationError{err: err}
}

func (c *ConnectionConfig) Validate() error {
	if c.Passphrase == "" {
		return NewValidationError(fmt.Errorf("passphrase is empty"))
	}
	if c.OrganizationId == "" {
		return NewValidationError(fmt.Errorf("organization id is empty"))
	}
	if c.CloudHost == "" {
		return NewValidationError(fmt.Errorf("cloud host is not set"))
	}
	if c.WSPort == 0 {
		return NewValidationError(fmt.Errorf("websocket port is not set"))
	}
	if c.WSPort < 0 {
		return NewValidationError(fmt.Errorf("invalid websocket port %d", c.WSPort))
	}
	return nil
}

func (c *ConnectionConfig) SetPassphrase(in string) {
	c.Passphrase = in
	config.Set("passphrase", in)
	c.save()
}
func (c *ConnectionConfig) SetOrganizationId(in string) {
	c.Passphrase = in
	config.Set("organization-id", in)
	c.save()
}

func (c *ConnectionConfig) save() {
	if err := config.WriteConfig(); err != nil {
		log.Error("Error while writing cloud config", "err", err, "config", c)
	}
}

func GetConnectionConfig() *ConnectionConfig {
	c := ConnectionConfig{
		Passphrase:     config.GetString("passphrase"),
		OrganizationId: config.GetString("organization-id"),
		CloudHost:      config.GetString("cloud-host"),
		WSPort:         config.GetInt("websocket-port"),
	}

	return &c
}

func GetReconnectTimeout() time.Duration {
	return config.GetDuration("reconnect-timeout")
}
