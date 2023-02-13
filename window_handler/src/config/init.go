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
}

var TargetSystemInfo = systemInfo{
	OS:              GET_INFO_FAILURE,
	SystemFramework: GET_INFO_FAILURE,
	MachineName:     GET_INFO_FAILURE,
}

func init() {
	loadInitConfigCache()
	_, err := os.Open(CONFOG_PATH)
	if err != nil {
		filePtr, _ := os.Create(CONFOG_PATH)
		addObserver()
		loadDefaultConfig()
		defer func() {
			filePtr.Close()
		}()
		return
	}
	loadConfig()
	addObserver()
}

// TODO 配置新增后的版本升级处理
func loadDefaultConfig() {
	defaultConfig := systemConfig{
		Version: version,
		QnqTarget: qnqTarget{
			Ip:        "0.0.0.0",
			LocalPath: "Not Set",
		},
		LocalSingleSync: localSync{
			SourcePath: "Not Set",
			TargetPath: "Not Set",
			Speed:      "Not Set",
			CheckMd5:   false,
		},
		LocalBatchSync: localSync{
			SourcePath: "Not Set",
			TargetPath: "Not Set",
			PeriodicSync: PeriodicSyncPolicy{
				Cycle:  time.Hour,
				Rate:   1,
				Enable: false,
			},
			TimingSync: TimingSyncPolicy{},
			Speed:      "Not Set",
			CheckMd5:   false,
		},
		VarianceAnalysis: varianceAnalysis{
			TimeStamp: true,
			Md5:       true,
		},
	}
	SystemConfigCache.Set(defaultConfig)
}
