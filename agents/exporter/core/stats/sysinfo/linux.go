package sysinfo

import (
	"cloudview/agents/exporter/core/iox"
	"errors"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// SysInfo represents the linux system info.
type SysInfo struct {
	Hostname   string    `json:"hostname"`
	FQDN       string    `json:"fqdn"`
	Domain     string    `json:"domain"`
	OsType     string    `json:"ostype"`
	OsRelease  string    `json:"osrelease"`
	OsVersion  string    `json:"osversion"`
	OsArch     string    `json:"osarch"`
	Uptime     float64   `json:"uptime"`
	InstanceId string    `json:"instanceId"`
	LocalIp    string    `json:"localIp"`
	Timestamp  time.Time `json:"timestamp"`
}

// getSysInfo gets the system info.
func GetSysInfo() (sysInfo SysInfo, err error) {
	sysInfo = SysInfo{}

	// Instance private ip
	sysInfo.LocalIp = getLocalIp()

	// Instance ID
	sysInfo.InstanceId = getInstanceId()
	// Hostname
	hostname, err := getHostname()
	if err != nil {
		return sysInfo, err
	}
	sysInfo.Hostname = hostname

	// Domain
	domain, err := getDomain()
	if err != nil {
		return sysInfo, err
	}
	sysInfo.Domain = domain

	// OS type
	osType, err := getOsType()
	if err != nil {
		return sysInfo, err
	}
	sysInfo.OsType = osType

	// OS relase
	osRelease, err := getOsRelease()
	if err != nil {
		return sysInfo, err
	}
	sysInfo.OsRelease = osRelease

	// OS version
	osVersion, err := getOsVersion()
	if err != nil {
		return sysInfo, err
	}
	sysInfo.OsVersion = osVersion

	// OS arch
	osArch, err := getOsArch()
	if err != nil {
		return sysInfo, err
	}
	sysInfo.OsArch = osArch

	// Uptime
	uptime, err := getUptime()
	if err != nil {
		return sysInfo, err
	}
	sysInfo.Uptime = uptime

	// FQDN
	fqdn, err := getFqdn()
	if err != nil {
		return sysInfo, err
	}
	sysInfo.FQDN = fqdn

	return sysInfo, nil
}

func getInstanceId() string {
	cg, err := currentProvider()
	if err != nil {
		return ""
	}
	return cg.instanceId()
}

func getHostname() (hostname string, err error) {
	return iox.ReadText("/proc/sys/kernel/hostname")
}

func getDomain() (domain string, err error) {
	return iox.ReadText("/proc/sys/kernel/domainname")
}

func getOsType() (osType string, err error) {
	return iox.ReadText("/proc/sys/kernel/ostype")
}

func getOsRelease() (osRelease string, err error) {
	return iox.ReadText("/proc/sys/kernel/osrelease")
}

func getOsVersion() (osVersion string, err error) {
	return iox.ReadText("/proc/sys/kernel/version")
}
func getLocalIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP.String()
}
func getOsArch() (osArch string, err error) {
	// Check `uname` path
	uname, err := exec.LookPath("uname")
	if err != nil {
		return "", err
	}

	// Run `uname -m` to get the OS architecture
	out, err := exec.Command(uname, "-m").Output()
	if err != nil {
		return "", err
	}

	osArch = strings.TrimSpace(string(out))
	return osArch, nil
}

func getUptime() (uptime float64, err error) {
	fields, err := iox.ReadTextLines("/proc/uptime")
	if err != nil {
		return -1, err
	}

	if len(fields) != 2 {
		return -1, errors.New("Error parsing /proc/uptime. It should have 2 fields")
	}

	uptime, err = strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return -1, err
	}

	return uptime, nil
}

func getFqdn() (fqdn string, err error) {
	// Check `hostname` path
	hostname, err := exec.LookPath("hostname")
	if err != nil {
		return "", err
	}

	// Run `hostname -f` to get the FQDN
	out, err := exec.Command(hostname, "-f").Output()
	if err != nil {
		return "", err
	}

	fqdn = strings.TrimSpace(string(out))
	return fqdn, nil
}
