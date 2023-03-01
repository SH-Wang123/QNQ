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
	QnqSTarget       qnqTarget        `json:"qnq_s_target"`
	QnqBTarget       qnqTarget        `json:"qnq_b_target"`
	LocalSingleSync  localSync        `json:"local_single_sync"`
	LocalBatchSync   localSync        `json:"local_batch_sync"`
	VarianceAnalysis varianceAnalysis `json:"variance_analysis"`
}

type localSync struct {
	SourcePath string     `json:"source_path"`
	TargetPath string     `json:"target_path"`
	SyncPolicy syncPolicy `json:"sync_policy"`
	Speed      string     `json:"speed"`
	CheckMd5   bool       `json:"check_md5"`
}

type syncPolicy struct {
	PolicySwitch bool               `json:"policy_switch"`
	PeriodicSync periodicSyncPolicy `json:"periodic_policy"`
	TimingSync   timingSyncPolicy   `json:"timing_sync"`
}

type periodicSyncPolicy struct {
	Cycle  time.Duration `json:"sync_cycle"`
	Rate   int           `json:"sync_rate"`
	Enable bool          `json:"enable"`
}

type timingSyncPolicy struct {
	Days   [7]bool `json:"sync_days"`
	Hour   uint8   `json:"hour"`
	Minute uint8   `json:"minute"`
	Enable bool    `json:"enable"`
}

type varianceAnalysis struct {
	TimeStamp bool `json:"time_stamp"`
	Md5       bool `json:"md5"`
}

type qnqTarget struct {
	Ip         string     `json:"qnq_target_ip"`
	LocalPath  string     `json:"local_path"`
	RemotePath string     `json:"remote_path"`
	SyncPolicy syncPolicy `json:"sync_policy"`
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

func (a *cacheConfig) GetSyncPolicy(isBatch bool, isRemote bool) *syncPolicy {
	if isBatch {
		if isRemote {
			return &a.Cache.QnqSTarget.SyncPolicy
		}
		return &a.Cache.LocalBatchSync.SyncPolicy
	} else {
		if isRemote {
			return &a.Cache.QnqBTarget.SyncPolicy
		}
		return &a.Cache.LocalSingleSync.SyncPolicy
	}
}

type LocalConfig struct {
}
