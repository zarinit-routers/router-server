package system

import (
	"errors"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
	"github.com/zarinit-routers/cli/df"
	"github.com/zarinit-routers/cli/nmcli"
	"github.com/zarinit-routers/router-server/pkg/cli/ip"
	"github.com/zarinit-routers/router-server/pkg/models"
)

func GetOSInfo(_ models.JSONMap) (any, error) {

	info, err := newInfo()
	if err != nil {
		return nil, err
	}
	return info, nil

}

type NetworkStats struct {
	network.Stats
	MAC string `json:"MAC"`
	IP  string `json:"IP"`
}

type OSInfo struct {
	Memory       *memory.Stats
	CpuStats     *cpu.Stats
	NetworkStats []NetworkStats
	DiskStats    []df.DiskStats
	LoadAverage  *loadavg.Stats
}

func newInfo() (OSInfo, error) {

	mem, memErr := memory.Get()
	cpuS, cpuErr := cpu.Get()
	net, netErr := network.Get()
	load, loadErr := loadavg.Get()
	if err := errors.Join(memErr, cpuErr, netErr, loadErr); err != nil {
		return OSInfo{}, err
	}
	return OSInfo{
		CpuStats:     cpuS,
		NetworkStats: toNetworkStatsArr(net),
		Memory:       mem,
		DiskStats:    df.Stats(),
		LoadAverage:  load}, nil
}

func toNetworkStatsArr(stats []network.Stats) []NetworkStats {
	stats2 := []NetworkStats{}
	for _, s := range stats {
		if s.Name == "" {
			continue
		}
		stats2 = append(stats2, newNetworkStats(s))
	}
	return stats2

}

func newNetworkStats(s network.Stats) NetworkStats {
	mac, _ := nmcli.GetHardwareAddress(s.Name)
	ipAddr, _ := ip.GetIP(s.Name)
	return NetworkStats{
		Stats: s,
		MAC:   mac,
		IP:    ipAddr,
	}
}
