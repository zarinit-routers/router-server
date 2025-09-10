package timezone

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/zarinit-routers/cli"
)

type TimedateInfoDictionary map[string]string

func getInfo() (TimedateInfoDictionary, error) {
	output, err := cli.Execute("timedatectl", "show")
	if err != nil {
		return nil, fmt.Errorf("failed to execute timedatectl show: %w", err)
	}
	info := TimedateInfoDictionary{}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		info[parts[0]] = parts[1]
	}
	return info, nil
}

func (t TimedateInfoDictionary) GetTimeZone() string {
	return t["Timezone"]
}

func (t TimedateInfoDictionary) NTP() bool {
	return strings.ToLower(t["NTP"]) == "yes"
}
