package config

type Config struct {
	Dhcp4 struct {
		InterfacesConfig struct {
			Interfaces []string `json:"interfaces"`
		} `json:"interfaces-config"`
	} `json:"Dhcp4"`
}
