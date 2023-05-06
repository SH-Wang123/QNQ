package config

import (
	"fmt"
	"net"
	"os"
)

func getLocalMachineName() string {
	machineName, _ := os.Hostname()
	return machineName
}

func getMac() (macs []string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, inter := range interfaces {
		if inter.HardwareAddr != nil {
			macs = append(macs, fmt.Sprintf("%v", inter.HardwareAddr))
		}
	}
	return macs
}

func getIp() (ips []string) {
	interfacesAddr, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range interfacesAddr {
		ipNet, isVailIpNet := address.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

func getTargetSystemInfo() string {
	machineName, _ := os.Hostname()
	return machineName
}

func loadInitConfigCache() {
	SystemConfigCache = cacheConfig{}

}

func addObserver() {
	SystemConfigCache.Register(&LocalConfigObserver{
		name: "local_system_config_observer",
	})
}

func GetCsvStr(str ...string) string {
	ret := ""
	for _, v := range str {
		if ret == "" {
			ret = v
		} else {
			ret += "," + v
		}
	}
	return ret
}
