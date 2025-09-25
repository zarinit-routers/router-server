package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

const (
	OptionKeyDriver              = "driver"
	OptionKeyCountryCode         = "country_code" // IDK is this even necessary to set
	OptionKeyWMMEnabled          = "wmm_enabled"
	OptionKeyInterface           = "interface"
	OptionKeyMACAddressACL       = "macaddr_acl"
	OptionKeyHTCapabilities      = "ht_capab" // A list of the 802.11n features supported by your device
	OptionKeyIgnoreBroadcastSSID = "ignore_broadcast_ssid"
	OptionKeySSID                = "ssid"
	OptionKeyAuthAlgorithms      = "auth_algs"
	OptionKeyIEEE80211d          = "ieee80211d"
	OptionKeyIEEE80211n          = "ieee80211n"
	OptionKeyChannel             = "channel"
	OptionKeyHWMode              = "hw_mode"
	OptionKeyWPA                 = "wpa"
	OptionKeyWPAPassphrase       = "wpa_passphrase"
	OptionKeyWPAKeyManagement    = "wpa_key_mgmt"
	OptionKeyRSNPairwise         = "rsn_pairwise"
	OptionKeyWPAPairwise         = "wpa_pairwise"
	OptionKeyBSS                 = "bss"
)

func (c *Config) GetSSID() string {
	return c.getOption(OptionKeySSID)
}

func (c *Config) SetSSID(ssid string) error {
	if err := containsSpecialRuns(ssid); err != nil {
		return fmt.Errorf("invalid SSID: %s", err)
	}
	c.setOption(OptionKeySSID, ssid)
	return nil
}

func (c *Config) GetInterface() string {
	return c.getOption(OptionKeyInterface)
}
func (c *Config) SetInterface(iface string) error {
	if err := containsSpecialRuns(iface); err != nil {
		return fmt.Errorf("invalid interface: %s", err)
	}
	if iface == "" {
		return fmt.Errorf("interface can't be empty")
	}
	c.setOption(OptionKeyInterface, iface)
	return nil
}

func containsSpecialRuns(s string) error {
	if strings.ContainsAny(s, "=") {
		return fmt.Errorf("string %q contains '=' this character not allowed for usage in this hostapd config parser/generator", s)
	}
	if strings.ContainsAny(s, "\t\n\r") {
		return fmt.Errorf("string %q contain not allowed whitespace characters", s)
	}
	return nil
}

type Driver string

const (
	DriverNL80211 Driver = "nl80211"

	DefaultDriver = DriverNL80211
)

func (c *Config) GetDriver() string {
	return c.getOption(OptionKeyDriver)
}

func (c *Config) SetDriver(driver Driver) error {
	if err := containsSpecialRuns(string(driver)); err != nil {
		return fmt.Errorf("invalid driver: %s", err)
	}
	c.setOption(OptionKeyDriver, string(driver))
	return nil
}

func (c *Config) GetChannel() int {
	ch, err := strconv.Atoi(c.getOption(OptionKeyChannel))
	if err != nil {
		log.Error("Invalid channel", "error", err, "channel", c.getOption(OptionKeyChannel))
		log.Warn("Using default channel", "channel", 0)
		return 0
	}
	return ch
}

const (
	Min2_4GChannel = 1
	Max2_4GChannel = 14
	Min5GChannel   = 1
	Max5GChannel   = 37
)

func (c *Config) SetChannel(channel int) error {
	switch c.GetHWMode() {
	case HWMode5G:
		if err := valueInBetween(channel, Min5GChannel, Max5GChannel); err != nil {
			return fmt.Errorf("invalid channel %d specified for 5G mode: %s", channel, err)
		}
	case HWMode2_4G:
		if err := valueInBetween(channel, Min2_4GChannel, Max2_4GChannel); err != nil {
			return fmt.Errorf("invalid channel %d specified for 2.4G mode: %s", channel, err)
		}
	default:
		return fmt.Errorf("unknown HWMode %s, can't validate chanel", c.GetHWMode())
	}
	c.setOption(OptionKeyChannel, fmt.Sprintf("%d", channel))
	return nil
}

func valueInBetween(val, min, max int) error {
	if val >= min && val <= max {
		return fmt.Errorf("value must be between %d and %d, got %d", min, max, val)
	}
	return nil
}

type HWMode string

const (
	HWMode5G   HWMode = "a"
	HWMode2_4G HWMode = "g"
)

func (c *Config) GetHWMode() HWMode {
	val := c.getOption(OptionKeyHWMode)
	return HWMode(val)
}

func (c *Config) SetHWMode(mode HWMode) {
	c.setOption(OptionKeyHWMode, string(mode))
	switch mode {
	case HWMode5G:
		c.SetChannel(Min5GChannel)
	case HWMode2_4G:
		c.SetChannel(Min2_4GChannel)
	}
}

type WPA string

const (
	WPA2 WPA = "2"
)

func (c *Config) SetWPA(w WPA) {
	c.setOption(OptionKeyWPA, string(w))
}

type WPAKeyManagement string

const (
	WPAKeyManagementPSK WPAKeyManagement = "WPA-PSK"
)

func (c *Config) SetWPAKeyManagement(m WPAKeyManagement) {
	c.setOption(OptionKeyWPAKeyManagement, string(m))
}

func (c *Config) GetWPAKeyManagement() WPAKeyManagement {
	return WPAKeyManagement(c.getOption(OptionKeyWPA))
}
func (c *Config) GetPassphrase() string {
	return c.getOption(OptionKeyWPAPassphrase)
}
func (c *Config) SetPassphrase(passphrase string) error {
	if err := containsSpecialRuns(passphrase); err != nil {
		return fmt.Errorf("invalid passphrase: %s", err)
	}
	if len(passphrase) < 8 {
		return fmt.Errorf("passphrase must be at least 8 characters long")
	}
	c.setOption(OptionKeyWPAPassphrase, passphrase)
	return nil
}
func (c *Config) GetSSIDVisibility() bool {
	return c.getOption(OptionKeyIgnoreBroadcastSSID) != "1"
}

func (c *Config) SetSSIDVisibility(visible bool) {
	if visible {
		c.setOption(OptionKeyIgnoreBroadcastSSID, "0")
	} else {
		c.setOption(OptionKeyIgnoreBroadcastSSID, "1")
	}
}

const (
	AuthAlgorithmWPA  = "1"
	AuthAlgorithmWEP  = "2"
	AuthAlgorithmBoth = "3"
)

func (c *Config) SetNecessaryDefaultOptions() error {
	c.SetDriver(DefaultDriver)
	if c.GetSSID() == "" {
		c.SetSSID("zarinit-wifi")
	}
	if err := c.SetInterface(viper.GetString("wifi-hotspot.interface")); err != nil {
		return fmt.Errorf("failed to set interface: %s", err)
	}
	c.SetWPA(WPA2)
	c.SetWPAKeyManagement(WPAKeyManagementPSK)

	if c.GetPassphrase() == "" {
		if err := c.SetPassphrase("1234567890"); err != nil {
			return fmt.Errorf("failed to set passphrase: %s", err)
		}
	}

	// switch c.getOption() {
	// case WIFI2_4G:
	// 	c.SetHWMode(HWMode2_4G)
	// case WIFI5G:
	// 	c.SetHWMode(HWMode5G)
	// default:
	// 	return fmt.Errorf("unknown WI-FI type")
	// }

	c.setOption(OptionKeyCountryCode, "RU")
	c.setOption(OptionKeyWMMEnabled, "1")
	c.setOption(OptionKeyRSNPairwise, "CCMP")
	c.setOption(OptionKeyWPAPairwise, "CCMP")
	c.setOption(OptionKeyAuthAlgorithms, AuthAlgorithmWPA)
	c.setOption(OptionKeyIEEE80211n, "1")
	c.setOption(OptionKeyIEEE80211d, "1")
	c.setOption(OptionKeyMACAddressACL, "0")
	return nil
}
