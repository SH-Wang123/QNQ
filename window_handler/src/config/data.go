package config

import (
	"time"
	"window_handler/common"
)

const GET_INFO_FAILURE = "Getting information failure"
const NOT_SET_STR = "Not Set"
const Version = "V0.0.4.5"
const WindowHeight = 600
const WindowWidth = 900

var SystemConfigCache cacheConfig

type systemConfig struct {
	Version          string           `json:"version"`
	QNQNetCells      []configNetCell  `json:"qnq_net_cells"`
	LocalSingleSync  localSync        `json:"local_single_sync"`
	LocalBatchSync   localSync        `json:"local_batch_sync"`
	PartitionSync    localSync        `json:"local_partition_sync"`
	CdpSnapshot      cdpSnapshot      `json:"cdp_snapshot"`
	VarianceAnalysis varianceAnalysis `json:"variance_analysis"`
	SystemSetting    systemSetting    `json:"system_setting"`
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

type configNetCell struct {
	Ip         string     `json:"qnq_target_ip"`
	SyncPolicy syncPolicy `json:"sync_policy"`
	Status     int        `json:"status"`
	Mark       string     `json:"mark"`
}

type systemInfo struct {
	OS              string
	SystemFramework string
	MachineName     string
	MAC             []string
	IP              []string
}

type cdpSnapshot struct {
	SourcePath string     `json:"cdp_source"`
	TargetPath string     `json:"cdp_target"`
	Policy     syncPolicy `json:"cdp_policy"`
}

type systemSetting struct {
	EnableOLog bool `json:"enable_o_log"`
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

func (a *cacheConfig) GetLocalSyncPolicy(isBatch bool, isPartitionSync bool) *syncPolicy {
	if isPartitionSync {
		return &a.Cache.PartitionSync.SyncPolicy
	} else if isBatch {
		return &a.Cache.LocalBatchSync.SyncPolicy
	} else {
		return &a.Cache.LocalSingleSync.SyncPolicy
	}
}

type LocalConfig struct {
}

func (c *configNetCell) GetTargetStatus() bool {
	return c.Status&10 >= 10
}

func (c *configNetCell) GetServerStatus() bool {
	return c.Status&01 == 1
}

func (s *systemConfig) AddNilNetCell() {
	cell := configNetCell{
		Ip:     "0.0.0.0",
		Status: 0,
		Mark:   "",
	}
	s.QNQNetCells = append(s.QNQNetCells, cell)
}
