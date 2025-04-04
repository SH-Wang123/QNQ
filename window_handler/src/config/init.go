package config

import (
	"os"
	"runtime"
	"time"
)

var LocalSystemInfo = systemInfo{
	OS:              runtime.GOOS,
	SystemFramework: runtime.GOARCH,
	MachineName:     getLocalMachineName(),
	MAC:             getMac(),
	IP:              getIp(),
}

var TargetSystemInfo = systemInfo{
	OS:              GET_INFO_FAILURE,
	SystemFramework: GET_INFO_FAILURE,
	MachineName:     GET_INFO_FAILURE,
}

func init() {
	loadInitConfigCache()
	addObserver()
	_, err := os.Open(CONFOG_PATH)
	if err != nil {
		filePtr, _ := os.Create(CONFOG_PATH)
		loadDefaultConfig()
		defer func() {
			filePtr.Close()
		}()
	} else {
		loadConfig()
	}
	initOLog()
	initQTP()
	addObserver()
	loadConfig()
}

// TODO 配置新增后的版本升级处理
func loadDefaultConfig() {
	defaultSyncPolicy := syncPolicy{
		PolicySwitch: false,
		PeriodicSync: periodicSyncPolicy{
			Cycle:  time.Hour,
			Rate:   1,
			Enable: false,
		},
		TimingSync: timingSyncPolicy{
			Days:   [7]bool{false, false, false, false, false, false, false},
			Hour:   15,
			Minute: 0,
			Enable: false,
		},
	}
	defaultLocal := localSync{
		SourcePath: NOT_SET_STR,
		TargetPath: NOT_SET_STR,
		SyncPolicy: defaultSyncPolicy,
		Speed:      NOT_SET_STR,
		CheckMd5:   false,
	}

	defaultConfig := systemConfig{
		Version:         Version,
		QNQNetCells:     make([]configNetCell, 0),
		LocalSingleSync: defaultLocal,
		LocalBatchSync:  defaultLocal,
		PartitionSync:   defaultLocal,
		VarianceAnalysis: varianceAnalysis{
			TimeStamp: true,
			Md5:       true,
		},
		SystemSetting: systemSetting{
			EnableOLog: true,
		},
	}
	SystemConfigCache.Set(defaultConfig)
}
