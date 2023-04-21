package config

import "os"

func getLocalMachineName() string {
	machineName, _ := os.Hostname()
	return machineName
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
