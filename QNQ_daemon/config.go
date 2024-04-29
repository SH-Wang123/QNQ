package QNQ

import (
	"encoding/json"
	"os"
	"time"
)

var ConfigCache systemConfig

const getInfoFailure = "Getting information failure"
const notSetStr = "Not Set"
const Version = "V0.0.1"
const configPath = "./qnq_config"

type observer interface {
	updateAd(any)
	getName() string
}

type localConfigObserver struct {
	name string
}

func (l *localConfigObserver) updateAd(a any) {
	config := a.(systemConfig)
	inputConfig(config)
}

func (l *localConfigObserver) getName() string {
	return l.name
}

type subject interface {
	register(observer)
	deregister(observer)
	notifyAll()
}

type systemConfig struct {
	observers       []observer
	Version         string        `json:"version"`
	LocalSingleSync localSync     `json:"local_single_sync"`
	LocalBatchSync  localSync     `json:"local_batch_sync"`
	RemoteQNQ       []target      `json:"remote_qnq"`
	SystemSetting   systemSetting `json:"system_setting"`
	sysInfo         sysInfo
}

func (a *systemConfig) register(o observer) {
	if o != nil {
		a.observers = append(a.observers, o)
	}
}

func (a *systemConfig) deregister(o observer) {

}

// notifyAll 通知观察者持久化配置
func (a *systemConfig) notifyAll() {
	for _, o := range a.observers {
		o.updateAd(*a)
	}
}

func (a *systemConfig) setLocalSingleSync(source, target string) {
	ConfigCache.LocalSingleSync.SourcePath = source
	ConfigCache.LocalSingleSync.TargetPath = target
	ConfigCache.notifyAll()
}

func (a *systemConfig) setLocalBatchSync(source, target string) {
	ConfigCache.LocalBatchSync.SourcePath = source
	ConfigCache.LocalBatchSync.TargetPath = target
	ConfigCache.notifyAll()
}

type localSync struct {
	SourcePath       string           `json:"source_path"`
	TargetPath       string           `json:"target_path"`
	SyncPolicy       syncPolicy       `json:"sync_policy"`
	Speed            string           `json:"speed"`
	VarianceAnalysis varianceAnalysis `json:"variance_analysis"`
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
	Days   uint8 `json:"sync_days"`
	Hour   uint8 `json:"hour"`
	Minute uint8 `json:"minute"`
	Enable bool  `json:"enable"`
}

type systemSetting struct {
	EnableOLog bool `json:"enable_o_log"`
}

type varianceAnalysis struct {
	TimeStamp bool `json:"time_stamp"`
	Md5       bool `json:"md5"`
}

type target struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	Des  string `json:"des"`
}

func initConfig() {
	loadConfig()
}

func registerConfigObserve() {
	localObserver := localConfigObserver{name: "localConfigObserver"}
	ConfigCache.register(&localObserver)
}

func inputConfig(config systemConfig) {
	result, _ := json.MarshalIndent(config, "", "    ")
	_ = os.WriteFile(configPath, result, 0644)
}

func readConfig() systemConfig {
	var config systemConfig
	bytes, _ := os.ReadFile(configPath)
	_ = json.Unmarshal(bytes, &config)
	return config
}

func loadConfig() {
	ConfigCache = readConfig()
	registerConfigObserve()
	if ConfigCache.Version != Version {
		ConfigCache.Version = Version
		ConfigCache.notifyAll()
	}
}
