package components

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type GPU struct {
	Vendor string // amd, intel, or nvidia
	Index  int    // GPU index for multi-GPU systems
}

func (g *GPU) detectVendor() error {
	// Check for AMD GPU
	amdPath := fmt.Sprintf("/sys/class/drm/card%d/device/", g.Index)
	if _, err := os.Stat(filepath.Join(amdPath, "gpu_busy_percent")); err == nil {
		g.Vendor = "amd"
		return nil
	}

	// Check for Intel GPU
	intelPath := fmt.Sprintf("/sys/class/drm/card%d/device/power/", g.Index)
	if _, err := os.Stat(filepath.Join(intelPath, "rc6_residency_ms")); err == nil {
		g.Vendor = "intel"
		return nil
	}

	// NVIDIA requires proprietary driver tools
	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		g.Vendor = "nvidia"
		return nil
	}

	return fmt.Errorf("no supported GPU detected")
}

func (g *GPU) readFileAsInt(path string) (int, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}

	return value, nil
}

// GPU Usage Percentage
func GpuPerc(_ string) string {
	gpu := &GPU{Index: 0}
	if err := gpu.detectVendor(); err != nil {
		return UnknownStr
	}

	var usage int
	var err error

	switch gpu.Vendor {
	case "amd":
		usage, err = gpu.readFileAsInt(
			fmt.Sprintf("/sys/class/drm/card%d/device/gpu_busy_percent", gpu.Index),
		)
	case "intel":
		// Calculate from frequency
		current, _ := gpu.readFileAsInt(
			fmt.Sprintf("/sys/class/drm/card%d/device/gt_act_freq_mhz", gpu.Index),
		)
		max, _ := gpu.readFileAsInt(
			fmt.Sprintf("/sys/class/drm/card%d/device/gt_RP0_freq_mhz", gpu.Index),
		)
		if max > 0 {
			usage = (current * 100) / max
		}
	case "nvidia":
		// Requires nvidia-smi parsing
		cmd := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits")
		output, err := cmd.CombinedOutput()
		if err == nil {
			usage, _ = strconv.Atoi(strings.TrimSpace(string(output)))
		}
	}

	if err != nil || usage < 0 || usage > 100 {
		return UnknownStr
	}

	return fmt.Sprintf("%d%%", usage)
}

// GPU Memory Percentage
func GpuMemPerc(_ string) string {
	gpu := &GPU{Index: 0}
	if err := gpu.detectVendor(); err != nil {
		return UnknownStr
	}

	var used, total int
	var err error

	switch gpu.Vendor {
	case "amd":
		total, err = gpu.readFileAsInt(
			fmt.Sprintf("/sys/class/drm/card%d/device/mem_info_vram_total", gpu.Index),
		)
		used, _ = gpu.readFileAsInt(
			fmt.Sprintf("/sys/class/drm/card%d/device/mem_info_vram_used", gpu.Index),
		)
	case "intel":
		// Intel reports memory in bytes
		total, err = gpu.readFileAsInt(
			fmt.Sprintf("/sys/class/drm/card%d/device/mem_info_vram_total", gpu.Index),
		)
		used, _ = gpu.readFileAsInt(
			fmt.Sprintf("/sys/class/drm/card%d/device/mem_info_vram_used", gpu.Index),
		)
		total /= 1024 * 1024 // Convert bytes to MB
		used /= 1024 * 1024
	case "nvidia":
		cmd := exec.Command("nvidia-smi", "--query-gpu=memory.used,memory.total", "--format=csv,noheader,nounits")
		output, err := cmd.CombinedOutput()
		if err == nil {
			parts := strings.Split(strings.TrimSpace(string(output)), ", ")
			if len(parts) == 2 {
				used, _ = strconv.Atoi(strings.TrimSuffix(parts[0], " MiB"))
				total, _ = strconv.Atoi(strings.TrimSuffix(parts[1], " MiB"))
			}
		}
	}

	if err != nil || total == 0 {
		return UnknownStr
	}

	return fmt.Sprintf("%.1f%%", float64(used)/float64(total)*100)
}

func (g *GPU) readSysfsInt(path string) (int, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func (g *GPU) findHwmonPath(pattern string) (string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return "", fmt.Errorf("no hwmon path found")
	}
	return matches[0], nil
}

func GpuTemp(_ string) string {
	gpu := &GPU{Index: 0}
	if err := gpu.detectVendor(); err != nil {
		return UnknownStr
	}

	var temp int
	var err error

	switch gpu.Vendor {
	case "amd", "intel":
		hwmonPath, err := gpu.findHwmonPath(
			fmt.Sprintf("/sys/class/drm/card%d/device/hwmon/hwmon*/temp1_input", gpu.Index),
		)
		if err == nil {
			temp, err = gpu.readSysfsInt(hwmonPath)
			if gpu.Vendor == "amd" {
				temp /= 1000 // Convert millidegrees to Celsius for AMD
			}
		}
	case "nvidia":
		cmd := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu", "--format=csv,noheader,nounits")
		output, err := cmd.Output()
		if err == nil {
			temp, _ = strconv.Atoi(strings.TrimSpace(string(output)))
		}
	}

	if err != nil || temp < 0 {
		return UnknownStr
	}

	return fmt.Sprintf("%d°C", temp)
}
