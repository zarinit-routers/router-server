package wifihotspot

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"github.com/zarinit-routers/cli/iw"
	"github.com/zarinit-routers/cli/nmcli"
	"github.com/zarinit-routers/cli/systemctl"
	"github.com/zarinit-routers/router-server/pkg/models"
)

func getConnectionName() (string, error) {
	name := viper.GetString("wifi-hotspot.connection-name")
	if name == "" {
		return "", fmt.Errorf("wifi hotspot connection name not configured (configuration key is 'wifi-hotspot.connection-name')")
	}
	return name, nil
}

// TODO: Rework this function
//
// Deprecated: this function should be reworked.
func getPassword() string {
	viper.SetDefault("wifi-hotspot.password", "12345678")
	return viper.GetString("wifi-hotspot.password")
}

func Enable(_ models.JSONMap) (any, error) {
	ifName := viper.GetString("wifi-hotspot.interface")
	if ifName == "" {
		return nil, fmt.Errorf("wifi hotspot interface not configured (configuration key is 'wifi-hotspot.interface')")
	}
	connName, err := getConnectionName()
	if err != nil {
		return nil, fmt.Errorf("failed get connection name: %s", err)
	}
	conn, err := nmcli.CreateWirelessConnection(ifName, connName, getPassword())
	if err != nil {
		return nil, fmt.Errorf("failed create wireless connection: %s", err)
	}

	ipAddr := "192.168.1.1/24"
	log.Warn("Setting IP address to constant, remove this behavior ASAP", "ip", ipAddr)
	err = errors.Join(
		conn.SetIP4Method(nmcli.ConnectionIP4MethodShared),
		conn.SetIP4Address(ipAddr),
		// conn.SetDNSAddresses([]string{"8.8.8.8", "8.8.4.4"}),
		// conn.SetDHCPRange(net.IPv4(192, 168, 1, 100), net.IPv4(192, 168, 1, 200)),
		// conn.SetDHCPLeaseTime(3600),
	)
	if err != nil {
		return nil, fmt.Errorf("failed configure connection: %s", err)
	}
	if err := conn.Up(); err != nil {
		return nil, fmt.Errorf("failed enable interface %q: %s", ifName, err)
	}

	err = errors.Join(systemctl.Enable("dhcpd"), systemctl.Enable("hostapd"))
	if err != nil {
		return nil, fmt.Errorf("failed enable services: %s", err)
	}
	return models.JSONMap{"enabled": true}, nil
}

func Disable(_ models.JSONMap) (any, error) {
	connName, err := getConnectionName()
	if err != nil {
		return nil, fmt.Errorf("failed get connection name: %s", err)
	}
	conn, err := nmcli.GetConnection(connName)
	if err != nil {
		return nil, fmt.Errorf("failed get connection %q: %s", conn.Name, err)
	}
	if err := conn.Down(); err != nil {
		return nil, fmt.Errorf("failed disable interface %q: %s", conn.Name, err)
	}
	return models.JSONMap{"enabled": false}, nil
}

func getConnection() (*nmcli.WirelessConnection, error) {
	connName, err := getConnectionName()
	if err != nil {
		return nil, fmt.Errorf("failed get connection name: %s", err)
	}
	conn, err := nmcli.GetConnection(connName)
	if err != nil {
		return nil, fmt.Errorf("failed get connection %q: %s", connName, err)
	}
	wr, err := conn.AsWireless()
	if err != nil {
		return nil, fmt.Errorf("failed get wireless connection %q: %s", connName, err)
	}
	return wr, nil
}
func GetStatus(_ models.JSONMap) (any, error) {
	conn, err := getConnection()
	if err != nil {
		return nil, fmt.Errorf("failed get connection information: %s", err)
	}

	return models.JSONMap{"enabled": conn.IsActive(), "ssid": conn.GetSSID(), "password": conn.GetPassword(), "band": conn.GetBand(), "channel": conn.GetChanel(), "hidden": conn.IsHidden()}, nil
}

func SetSSID(args models.JSONMap) (any, error) {
	ssid, ok := args["ssid"].(string)
	if !ok {
		return nil, fmt.Errorf("ssid not specified")
	}
	conn, err := getConnection()
	if err != nil {
		return nil, fmt.Errorf("failed get connection information: %s", err)
	}

	err = conn.SetSSID(ssid)
	if err != nil {
		return nil, fmt.Errorf("failed set ssid: %s", err)
	}
	return models.JSONMap{"ssid": conn.GetSSID()}, nil
}

func SetSSIDVisibility(args models.JSONMap) (any, error) {
	hidden, ok := args["hidden"].(bool)
	if !ok {
		return nil, fmt.Errorf("visibility not specified (key 'hidden' of type bool)")
	}
	conn, err := getConnection()
	if err != nil {
		return nil, fmt.Errorf("failed get connection information: %s", err)
	}

	err = conn.SetHidden(hidden)
	if err != nil {
		return nil, fmt.Errorf("failed set ssid visibility: %s", err)
	}
	return models.JSONMap{"hidden": hidden}, nil
}

func SetPassword(args models.JSONMap) (any, error) {
	password, ok := args["password"].(string)
	if !ok {
		return nil, fmt.Errorf("password not specified")
	}
	conn, err := getConnection()
	if err != nil {
		return nil, fmt.Errorf("failed get connection information: %s", err)
	}

	err = conn.SetPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed set password: %s", err)
	}
	return models.JSONMap{"password": password}, nil
}

func SetChannel(args models.JSONMap) (any, error) {
	channel, ok := args["channel"].(int)
	if !ok {
		return nil, fmt.Errorf("channel not specified")
	}
	conn, err := getConnection()
	if err != nil {
		return nil, fmt.Errorf("failed get connection information: %s", err)
	}

	err = conn.SetChannel(channel)
	if err != nil {
		return nil, fmt.Errorf("failed set channel: %s", err)
	}
	return models.JSONMap{"channel": channel}, nil
}
func GetConnectedClients(_ models.JSONMap) (any, error) {
	conn, err := getConnection()
	if err != nil {
		return nil, fmt.Errorf("failed get connection information: %s", err)
	}

	clients, err := iw.GetConnectedDevices(conn.Device)
	if err != nil {
		return nil, fmt.Errorf("failed get connected clients: %s", err)
	}
	return models.JSONMap{"clients": clients}, nil
}
