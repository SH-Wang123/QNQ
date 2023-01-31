package config

import (
	"runtime"
	"window_handler/common"
)

const GET_INFO_FAILURE = "Getting information failure"
const NOT_SET_TARGET = "Not set QNQ target"

var SystemConfigCache cacheConfig

type systemConfig struct {
	QnqTarget        qnqTarget        `json:"qnq_target"`
	LocalSingleSync  localSync        `json:"local_single_sync"`
	LocalBatchSync   localSync        `json:"local_batch_sync"`
	VarianceAnalysis varianceAnalysis `json:"variance_analysis"`
}

type localSync struct {
	SourcePath string `json:"source_path"`
	TargetPath string `json:"target_path"`
	Speed      string `json:"speed"`
	CheckMd5   bool   `json:"check_md5"`
}

type varianceAnalysis struct {
	TimeStamp bool `json:"time_stamp"`
	FileSize  bool `json:"file_size"`
	LastTime  bool `json:"last_time"`
	Md5       bool `json:"md5"`
}

type qnqTarget struct {
	Ip        string `json:"qnq_target_ip"`
	LocalPath string `json:"local_path"`
}

type systemInfo struct {
	OS              string
	SystemFramework string
	MachineName     string
}

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

// CacheConfig Subject Object
type cacheConfig struct {
	Cache     systemConfig
	observers []common.Observer
}

func (a *cacheConfig) Register(observer common.Observer) {
	a.observers = append(a.observers, observer)
}

func (a *cacheConfig) Deregister(observer common.Observer) {

}

func (a *cacheConfig) NotifyAll() {
	for _, observer := range a.observers {
		observer.UpdateAd(a.Cache)
	}
}

func (a *cacheConfig) Value() systemConfig {
	return a.Cache
}

func (a *cacheConfig) Set(s systemConfig) {
	a.Cache = s
	a.NotifyAll()
}

type LocalConfig struct {
}
