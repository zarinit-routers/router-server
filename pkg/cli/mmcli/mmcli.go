package mmcli

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cli"
)

func Get(modem string) (*ModemInfo, error) {
	output, err := cli.Execute("mmcli", ModemFlag(modem), JsonOutputFlag)
	if err != nil {
		log.Error("Failed get modem info", "error", err, "modem", modem)
		return nil, err
	}
	info := struct {
		Modem ModemInfo `json:"modem"`
	}{}
	err = json.Unmarshal(output, &info)
	if err != nil {
		log.Error("Failed parse modem info from JSON", "error", err)
		return nil, err
	}
	return &info.Modem, nil
}

func (m *ModemInfo) Disable() error {
	_, err := cli.Execute("mmcli", ModemFlag(m.DBusPath), "--disable")
	return err
}

func (m *ModemInfo) Enable() error {
	m.SetPowerStateOn() // ensures that it turned on
	_, err := cli.Execute("mmcli", ModemFlag(m.DBusPath), "--enable")
	return err
}

func (m *ModemInfo) SetPowerStateOff() error {
	log.Debug("Set power state off", "modem", m.DBusPath)
	_, err := cli.Execute("mmcli", ModemFlag(m.DBusPath), "--set-power-state-off")
	return err
}

func (m *ModemInfo) SetPowerStateOn() error {
	log.Debug("Set power state on", "modem", m.DBusPath)
	_, err := cli.Execute("mmcli", ModemFlag(m.DBusPath), "--set-power-state-on")
	return err
}

func list() ([]string, error) {
	output, err := cli.Execute("mmcli", ListModemsFlag, JsonOutputFlag)
	if err != nil {
		log.Errorf("Failed get modems list: %s", err)
		return nil, err
	}
	modems := struct {
		List []string `json:"modem-list"`
	}{}
	err = json.Unmarshal(output, &modems)
	if err != nil {
		log.Errorf("Failed parse modem list from JSON: %s", err)
		return nil, err
	}
	return modems.List, nil
}
func List() ([]*ModemInfo, error) {
	modems, err := list()
	if err != nil {
		return nil, fmt.Errorf("failed get list of modems: %s", err)
	}
	var list []*ModemInfo
	for _, modem := range modems {
		info, err := Get(modem)
		if err != nil {
			log.Errorf("Failed get modem info: %s", err)
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (m *ModemInfo) GetBearer() (*BearerInfo, error) {
	if len(m.Generic.Bearers) == 0 {
		return nil, fmt.Errorf("no bearers")
	}
	_, err := cli.Execute("mmcli", BearerFlag(m.Generic.Bearers[0]), JsonOutputFlag)
	var info struct {
		Bearer BearerInfo `json:"bearer"`
	}
	log.Warn("GetBearer not implemented yet")
	return &info.Bearer, err
}
