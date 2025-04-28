package components

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RAM struct {
	MemTotal     uint64
	MemFree      uint64
	MemAvailable uint64
	SwapTotal    uint64
	SwapFree     uint64
}

func (ram *RAM) readRamInfo() error {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return fmt.Errorf("failed to open meminfo: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	targets := map[string]*uint64{
		"MemTotal":     &ram.MemTotal,
		"MemFree":      &ram.MemFree,
		"MemAvailable": &ram.MemAvailable,
		"SwapTotal":    &ram.SwapTotal,
		"SwapFree":     &ram.SwapFree,
	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		field, exists := targets[key]
		if !exists {
			continue
		}

		valueStr := strings.Fields(strings.TrimSpace(parts[1]))[0]
		value, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			return fmt.Errorf("parse error for %s: %w", key, err)
		}

		*field = value
	}

	if ram.MemTotal == 0 {
		return fmt.Errorf("missing required memory fields")
	}

	return nil
}

func (ram *RAM) Render() string {
	if err := ram.readRamInfo(); err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}

	return fmt.Sprintf(`MemTotal:       %7d kB
MemFree:        %7d kB
MemAvailable:   %7d kB
SwapTotal:      %7d kB
SwapFree:       %7d kB`,
		ram.MemTotal,
		ram.MemFree,
		ram.MemAvailable,
		ram.SwapTotal,
		ram.SwapFree,
	)
}

func RamFree(_ string) string {
	ram := &RAM{}
	if err := ram.readRamInfo(); err != nil {
		return UnknownStr
	}
	return fmt.Sprintf("%.1fG", float64(ram.MemFree)/1024/1024)
}

func RamPerc(_ string) string {
	ram := &RAM{}
	if err := ram.readRamInfo(); err != nil {
		return UnknownStr
	}
	return fmt.Sprintf("%.1f%%", float64(ram.MemTotal-ram.MemAvailable)/float64(ram.MemTotal)*100)
}
