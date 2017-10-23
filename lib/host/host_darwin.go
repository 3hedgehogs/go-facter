package host

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	h "github.com/shirou/gopsutil/host"
)

// Facter interface
type Facter interface {
	Add(string, interface{})
}

// capitalize the first letter of given string
func capitalize(label string) string {
	firstLetter := strings.SplitN(label, "", 2)
	if len(firstLetter) < 1 {
		return label
	}
	return fmt.Sprintf("%v%v", strings.ToUpper(firstLetter[0]),
		strings.TrimPrefix(label, firstLetter[0]))
}

// getUniqueID returns executes % hostid; and returns its STDOUT as a string.
func getUniqueID() (string, error) {
	cmd := exec.Command("/usr/bin/hostid")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(out.String(), "\n"), nil
}

// getMacOSKernelVersion returns Darwin kernel version.
func getMacOSKernelVersion() (string, error) {
	systemProfiler := exec.Command("/usr/sbin/system_profiler","SPSoftwareDataType")
	var out bytes.Buffer
	systemProfiler.Stdout = &out
	err := systemProfiler.Run()
	if err != nil {
		return "", err
	}
	f := func(c rune) bool {
		return c == ':'
	}
	for _, line := range strings.Split(out.String(), "\n") {
		values := strings.FieldsFunc(strings.TrimSpace(line), f)
		if len(values) != 2 {
			continue
		}
		if strings.HasPrefix(values[0], "Kernel Version") {
			return strings.TrimSpace(values[1]), nil
		}
	}
	return "", nil
}

// guessArch tries to guess architecture based on HW model
func guessArch(HWModel string) string {
	var arch string
	switch HWModel {
	case "x86_64":
		arch = "amd64"
		break
	default:
		arch = "unknown"
		break
	}
	return arch
}

// int8ToString converts [65]int8 in syscall.Utsname to string
func int8ToString(bs [65]int8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		if v < 0 {
			b[i] = byte(256 + int(v))
		} else {
			b[i] = byte(v)
		}
	}
	return strings.TrimRight(string(b), "\x00")
}

// GetHostFacts gathers facts related to Host
func GetHostFacts(f Facter) error {
	hostInfo, err := h.Info()
	if err != nil {
		return err
	}

	// os.Hostname is not fqdn, fix
	hostInfo.Hostname = GetFQDN()

	f.Add("fqdn", hostInfo.Hostname)
	splitted := strings.SplitN(hostInfo.Hostname, ".", 2)
	var hostname *string
	if len(splitted) > 1 {
		hostname = &splitted[0]
		f.Add("domain", splitted[1])
	} else {
		hostname = &hostInfo.Hostname
	}
	f.Add("hostname", *hostname)

	isVirtual := false
	if hostInfo.VirtualizationRole == "host" {
		isVirtual = false
	} else if hostInfo.VirtualizationRole != "" {
		isVirtual = true
	}
	if hostInfo.VirtualizationSystem == "" && !isVirtual {
		hostInfo.VirtualizationSystem = "physical"
	}
	f.Add("is_virtual", isVirtual)

	f.Add("kernel", capitalize(hostInfo.OS))
	f.Add("operatingsystemrelease", hostInfo.PlatformVersion)
	f.Add("operatingsystem", capitalize(hostInfo.Platform))
	if hostInfo.PlatformFamily == "" {
		hostInfo.PlatformFamily = hostInfo.OS
	}
	f.Add("osfamily", capitalize(hostInfo.PlatformFamily))
	f.Add("uptime_seconds", hostInfo.Uptime)
	f.Add("uptime_minutes", hostInfo.Uptime/60)
	f.Add("uptime_hours", hostInfo.Uptime/60/60)
	f.Add("uptime_days", hostInfo.Uptime/60/60/24)
	f.Add("uptime", fmt.Sprintf("%d days", hostInfo.Uptime/60/60/24))
	f.Add("virtual", hostInfo.VirtualizationSystem)

	envPath := os.Getenv("PATH")
	if envPath != "" {
		f.Add("path", envPath)
	}

	user, err := user.Current()
	if err == nil {
		f.Add("id", user.Username)
	} else {
		panic(err)
	}

	z, _ := time.Now().Zone()
	f.Add("timezone", z)

	hostid, err := getUniqueID()
	if err == nil {
		f.Add("uniqueid", hostid)
	}
	kernelVersion, err := getMacOSKernelVersion()
	if err == nil && kernelVersion != "" {
		f.Add("kernelversion", strings.Split(kernelVersion, " ")[1])
	}
	return nil
}
