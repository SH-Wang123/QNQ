package config

import "os"

func InitConfig() {
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
			Speed:      "Not Set",
			CheckMd5:   false,
		},
		VarianceAnalysis: varianceAnalysis{
			TimeStamp: true,
			FileSize:  true,
			LastTime:  true,
			Md5:       true,
		},
	}
	SystemConfigCache.Set(defaultConfig)
}
