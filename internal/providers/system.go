package providers

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/MartialM1nd/opnlab/internal/providers"
)

// SystemProvider gathers system metrics from FreeBSD.
type SystemProvider struct {
	*providers.BaseProvider
}

// NewSystemProvider creates a new system provider.
func NewSystemProvider() *SystemProvider {
	return &SystemProvider{
		BaseProvider: providers.NewBaseProvider("system"),
	}
}

// Collect gathers system metrics.
func (p *SystemProvider) Collect() (map[string]interface{}, error) {
	data := make(map[string]interface{})

	// Get CPU usage
	cpu, err := getCPUUsage()
	if err == nil {
		data["cpu"] = cpu
	}

	// Get memory usage
	mem, err := getMemoryUsage()
	if err == nil {
		data["memory"] = mem
	}

	// Get load average
	load, err := getLoadAverage()
	if err == nil {
		data["load"] = load
	}

	// Get disk usage
	disk, err := getDiskUsage()
	if err == nil {
		data["disk"] = disk
	}

	// Get uptime
	uptime, err := getUptime()
	if err == nil {
		data["uptime"] = uptime
	}

	data["timestamp"] = time.Now().Format(time.RFC3339)
	return data, nil
}

// getCPUUsage returns CPU usage percentage.
func getCPUUsage() (map[string]interface{}, error) {
	out, err := exec.Command("vmstat", "-c", "1").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "
")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected vmstat output")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 18 {
		return nil, fmt.Errorf("not enough fields in vmstat output")
	}

	user, _ := strconv.ParseFloat(fields[13], 64)
	system, _ := strconv.ParseFloat(fields[14], 64)
	idle, _ := strconv.ParseFloat(fields[16], 64)

	used := 100 - idle

	return map[string]interface{}{
		"user":   user,
		"system": system,
		"idle":   idle,
		"used":   used,
	}, nil
}

// getMemoryUsage returns memory usage.
func getMemoryUsage() (map[string]interface{}, error) {
	out, err := exec.Command("sysctl", "-n", 
		"hw.physmem",
		"vm.stats.vm.v_wire_count",
		"vm.stats.vm.v_active_count",
		"vm.stats.vm.v_inactive_count",
		"vm.stats.vm.v_cache_count",
		"vm.stats.vm.v_free_count",
	).Output()
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(out))
	if len(fields) < 6 {
		return nil, fmt.Errorf("not enough sysctl fields")
	}

	total, _ := strconv.ParseUint(fields[0], 10, 64)
	wire, _ := strconv.ParseUint(fields[1], 10, 64)
	active, _ := strconv.ParseUint(fields[2], 10, 64)
	inactive, _ := strconv.ParseUint(fields[3], 10, 64)
	cache, _ := strconv.ParseUint(fields[4], 10, 64)
	free, _ := strconv.ParseUint(fields[5], 10, 64)

	pageSize := uint64(4096)
	
	return map[string]interface{}{
		"total":     total,
		"wire":      wire * pageSize,
		"active":    active * pageSize,
		"inactive":  inactive * pageSize,
		"cache":     cache * pageSize,
		"free":      free * pageSize,
		"used":      (total - free*pageSize),
		"usedPercent": float64(total-free*pageSize) / float64(total) * 100,
	}, nil
}

// getLoadAverage returns load average.
func getLoadAverage() (map[string]interface{}, error) {
	out, err := exec.Command("sysctl", "-n", "vm.loadavg").Output()
	if err != nil {
		return nil, err
	}

	trimmed := strings.Trim(string(out), "{}
 ")
	parts := strings.Fields(trimmed)

	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected loadavg format")
	}

	load1, _ := strconv.ParseFloat(parts[0], 64)
	load5, _ := strconv.ParseFloat(parms[1], 64)
	load15, _ := strconv.ParseFloat(parts[2], 64)

	return map[string]interface{}{
		"1min":  load1,
		"5min":  load5,
		"15min": load15,
	}, nil
}

// getDiskUsage returns disk usage.
func getDiskUsage() ([]map[string]interface{}, error) {
	out, err := exec.Command("df", "-k").Output()
	if err != nil {
		return nil, err
	}

	var disks []map[string]interface{}
	lines := strings.Split(string(out), "
")

	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		if fields[0] == "tmpfs" || fields[0] == "devfs" || fields[0] == "fdescfs" {
			continue
		}

		mountPoint := fields[5]
		used, _ := strconv.ParseUint(fields[2], 10, 1024)
		available, _ := strconv.ParseUint(fields[3], 10, 1024)
		capacity, _ := strconv.Atoi(strings.Trim(fields[4], "%"))

		disks = append(disks, map[string]interface{}{
			"filesystem": fields[0],
			"mount":      mountPoint,
			"used":       used * 1024,
			"available":  available * 1024,
			"capacity":   capacity,
		})
	}

	return disks, nil
}

// getUptime returns system uptime.
func getUptime() (map[string]interface{}, error) {
	out, err := exec.Command("sysctl", "-n", "kern.boottime").Output()
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(out))
	bootTime, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return nil, err
	}

	boot := time.Unix(bootTime, 0)
	uptime := time.Since(boot)

	return map[string]interface{}{
		"bootTime":      boot.Format(time.RFC3339),
		"uptimeSeconds": uptime.Seconds(),
		"uptime":        uptime.String(),
	}, nil
}

// Actions returns available actions (none for system provider).
func (p *SystemProvider) Actions() map[string]providers.Action {
	return map[string]providers.Action{}
}

// Ensure Provider implements providers.Provider
var _ providers.Provider = (*SystemProvider)(nil)
