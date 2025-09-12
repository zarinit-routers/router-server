package mmcli

import (
	"encoding/json"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cli"
)

func (m *ModemInfo) GetSignal() (*ModemSignal, error) {
	output, err := cli.ExecuteWrap("mmcli", ModemFlag(m.DBusPath), "--signal-get", JsonOutputFlag)
	if err != nil {
		log.Warn("Error occurred while getting signal, and was ignored", "error", err)
		log.Warn("This error should not be ignored! Call Katy248 to fix it ASAP")
		// log.Errorf("Failed get signal: %s", err)
		// return nil, err
	}
	info := struct {
		Modem struct {
			Signal ModemSignal `json:"signal"`
		} `json:"modem"`
	}{}
	err = json.Unmarshal(output, &info)
	if err != nil {
		log.Errorf("Failed parse modem signal info from JSON: %s", err)
		return nil, err
	}
	log.Debug("Modem signal", "JSONString", string(output), "signal", info.Modem.Signal)
	return &info.Modem.Signal, nil
}
