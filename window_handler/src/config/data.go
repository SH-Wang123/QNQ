package config

import (
	"time"
	"window_handler/common"
)

const GET_INFO_FAILURE = "Getting information failure"
const NOT_SET_STR = "Not Set"
const version = "V0.0.3"
const WindowHeight = 600
const WindowWidth = 600

var SystemConfigCache cacheConfig

type systemConfig struct {
	Version          string           `json:"version"`
	QnqTarget        qnqTarget        `json:"qnq_target"`
	LocalSingleSync  localSync        `json:"local_single_sync"`
	LocalBatchSync   localSync        `json:"local_batch_sync"`
	VarianceAnalysis varianceAnalysis `json:"variance_analysis"`
}

type localSync struct {
	SourcePath   string             `json:"source_path"`
	TargetPath   string             `json:"target_path"`
	PolicySwitch bool               `json:"policy_switch"`
	PeriodicSync PeriodicSyncPolicy `json:"periodic_policy"`
	TimingSync   TimingSyncPolicy   `json:"timing_sync"`
	Speed        string             `json:"speed"`
	CheckMd5     bool               `json:"check_md5"`
}

type PeriodicSyncPolicy struct {
	Cycle  time.Duration `json:"sync_cycle"`
	Rate   int           `json:"sync_rate"`
	Enable bool          `json:"enable"`
}

type TimingSyncPolicy struct {
	Days   [7]bool   `json:"sync_days"`
	Time   time.Time `json:"sync_time"`
	Enable bool      `json:"enable"`
}

type varianceAnalysis struct {
	TimeStamp bool `json:"time_stamp"`
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

func (a *cacheConfig) GetLocalPeriodicSyncPolicy(isBatch bool) *localSync {
	if isBatch {
		return &a.Cache.LocalBatchSync
	} else {
		return &a.Cache.LocalSingleSync
	}
}

type LocalConfig struct {
}
