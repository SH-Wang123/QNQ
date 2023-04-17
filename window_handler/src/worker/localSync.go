package worker

import (
	"io"
	"log"
	"os"
	"sync"
	"time"
	"window_handler/common"
	"window_handler/config"
)

// TODO 及时释放这些引用
var totalSizeMap = make(map[string]uint64)
var doneSizeMap = make(map[string]uint64)
var batchErrMap = make(map[string][]string)
var resourceLock = &sync.Mutex{}

var syncErrorMap = make(map[string][]SyncFileError, 0)
var errorMapLock = &sync.Mutex{}
var (
	periodicLocalBatchTicker   *time.Ticker
	periodicLocalSingleTicker  *time.Ticker
	periodicPartitionTicker    *time.Ticker
	periodicRemoteBatchTicker  *time.Ticker
	periodicRemoteSingleTicker *time.Ticker
	timingLocalBatchTicker     *time.Ticker
	timingLocalSingleTicker    *time.Ticker
	timingPartitionTicker      *time.Ticker
	timingRemoteBatchTicker    *time.Ticker
	timingRemoteSingleTicker   *time.Ticker
)

var (
	COMPARE_RUNNING = "[Comparing] "
	SYNC_RUNNING    = "[Syncing] "
	VERIFY_MD5      = "[Verify MD5] "
)

func NewLocalSingleWorker(sourceFile *os.File, targetFile *os.File, sn string, md5CacheFlag bool) *common.QWorker {
	return &common.QWorker{
		SN:              sn,
		Sub:             nil,
		ExecuteFunc:     localSyncSingleFile,
		DeconstructFunc: closeAndCheckFile,
		PrivateFile:     sourceFile,
		TargetFile:      targetFile,
		Md5CacheFlag:    md5CacheFlag,
	}
}

func FileNode2TreeMap(data *map[string][]string) {
	datat := *data
	datat[""] = append(datat[""], LocalBSFileNode.AbstractPath)
	GetFileTreeMap(LocalBSFileNode, data)
}

func GetFileTreeMap(node *FileNode, data *map[string][]string) {
	datat := *data
	if !node.IsDirectory {
		return
	}
	for _, child := range node.ChildrenNodeList {
		var key string
		if node.VarianceType == VARIANCE_ROOT {
			key = node.AbstractPath
		} else {
			key = node.AnchorPointPath
		}
		datat[key] = append(datat[key], child.AnchorPointPath)
		GetFileTreeMap(child, data)
	}
}

func batchSyncFile(startPath string, targetPath string, sn *string, busType int, loadMd5Cache bool) {
	sf, err1 := OpenDir(startPath)
	if err1 != nil {
		return
	}
	sfInfo, _ := sf.Stat()
	sfMode := sfInfo.Mode()
	common.CreateDir(targetPath, &sfMode)
	tf, err2 := OpenDir(targetPath)
	if err1 != nil || err2 != nil {
		return
	}
	children, _ := sf.Readdir(-1)
	common.CloseFile(sf, tf)
	ReverseCompareAndDelete(startPath, targetPath)
	if !common.GetRunningFlag(busType) {
		cancelTask(busType, common.GetForceDoneSignal(busType))
		return
	}
	for _, child := range children {
		targetAbsPath := targetPath + fileSeparator + child.Name()
		sourceAbsPath := startPath + fileSeparator + child.Name()
		if !child.IsDir() {
			if child.Size() <= 512*int64(MB) {
				isExist, _ := common.IsExist(targetAbsPath)
				tf, err3 := common.OpenFile(targetAbsPath, true)
				rsf, err4 := common.OpenFile(sourceAbsPath, false)
				if err3 != nil || err4 != nil {
					log.Printf("sf err: %v, tf err: %v", err4, err3)
					continue
				}
				common.SetCurrentSyncFile(*sn, COMPARE_RUNNING, sourceAbsPath)
				if isExist && CompareMd5(tf, rsf) {
					fInfo, _ := tf.Stat()
					addSizeToDoneMap(*sn, uint64(fInfo.Size()))
					common.CloseFile(tf, rsf)
					continue
				}
				//reopen file
				common.CloseFile(tf, rsf)
			}
			tf, errT := common.OpenFile(targetAbsPath, true)
			rsf, errS := common.OpenFile(sourceAbsPath, false)
			//common.GetCoroutinesPool().Submit(worker.Execute())
			if errT == nil && errS == nil {
				worker := NewLocalSingleWorker(rsf, tf, *sn, loadMd5Cache)
				worker.Execute()
			} else {
				common.CloseFile(rsf, tf)
			}
		} else {
			batchSyncFile(sourceAbsPath, targetAbsPath, sn, busType, loadMd5Cache)
		}
	}
}

func preSyncSingleTime(busType int, runningTag int) (sn string, lock *sync.WaitGroup, startTime string) {
	sn = common.GetSNCount()
	lock = common.GetStartLock(busType)
	lock.Add(1)
	common.SetCurrentSN(busType, sn)
	common.SetRunningFlag(busType, true)
	common.SendSignal2GWChannel(runningTag)
	initSizeMap(sn)
	return sn, lock, getNowTimeStr()
}

func afterSyncSingleTime(busType int, doneTag int) {
	common.SetRunningFlag(busType, false)
	common.SendSignal2GWChannel(doneTag)
}

// LocalBatchSyncSingleTime 直接读取配置文件，无需参数
func LocalBatchSyncSingleTime(isPolicy bool) {
	if isPolicy {
		if common.GetRunningFlag(common.TYPE_LOCAL_BATCH) {
			return
		}
	}
	sn, lock, startTime := preSyncSingleTime(common.TYPE_LOCAL_BATCH, common.LOCAL_BATCH_RUNNING)
	defer afterSyncSingleTime(common.TYPE_LOCAL_BATCH, common.LOCAL_BATCH_FORCE_DONE)
	GetTotalSize(&sn, config.SystemConfigCache.Cache.LocalBatchSync.SourcePath, true, lock)
	lock.Done()
	batchSyncFile(
		config.SystemConfigCache.Cache.LocalBatchSync.SourcePath,
		config.SystemConfigCache.Cache.LocalBatchSync.TargetPath,
		&sn,
		common.TYPE_LOCAL_BATCH,
		false,
	)
	recordLog(
		common.TYPE_LOCAL_BATCH,
		startTime,
		config.SystemConfigCache.Cache.LocalBatchSync.TargetPath,
		config.SystemConfigCache.Cache.LocalBatchSync.SourcePath)
}

func cancelTask(busType int, doneTag int) {
	afterSyncSingleTime(busType, doneTag)
}

// CancelTask 由外界强制设置任务终止标志
func CancelTask(busType int) {
	common.SetRunningFlag(busType, false)
}

// LocalSingleSyncSingleTime 直接读取配置文件，无需参数
func LocalSingleSyncSingleTime(isPolicy bool) {
	if isPolicy {
		if common.GetRunningFlag(common.TYPE_LOCAL_SING) {
			return
		}
	}
	sn, lock, startTime := preSyncSingleTime(common.TYPE_LOCAL_SING, common.LOCAL_SINGLE_RUNNING)
	defer afterSyncSingleTime(common.TYPE_LOCAL_SING, common.LOCAL_SINGLE_FORCE_DONE)
	sf, _ := common.OpenFile(config.SystemConfigCache.Cache.LocalSingleSync.SourcePath, false)
	tf := getSingleTargetFile(sf, config.SystemConfigCache.Cache.LocalSingleSync.TargetPath)
	sfInfo, err := sf.Stat()
	if err == nil {
		addSizeToTotalMap(sn, uint64(sfInfo.Size()))
	} else {
		log.Printf("single sync get file stat err: %v", err)
	}
	lock.Done()
	worker := NewLocalSingleWorker(sf, tf, sn, false)
	worker.Execute()
	recordLog(
		common.TYPE_LOCAL_SING,
		startTime,
		tf.Name(),
		sf.Name(),
	)
}

func periodicLocalBatchSync() {
	if config.SystemConfigCache.Cache.LocalBatchSync.SyncPolicy.PeriodicSync.Enable {
		if config.SystemConfigCache.Value().LocalBatchSync.SyncPolicy.PolicySwitch {
			LocalBatchSyncSingleTime(true)
		}
	}
}

func periodicLocalSingleSync() {
	if config.SystemConfigCache.Cache.LocalSingleSync.SyncPolicy.PeriodicSync.Enable {
		if config.SystemConfigCache.Cache.LocalSingleSync.SyncPolicy.PolicySwitch {
			LocalSingleSyncSingleTime(true)
		}
	}
}

func periodicPartitionSync() {
	if config.SystemConfigCache.Cache.PartitionSync.SyncPolicy.PeriodicSync.Enable {
		if config.SystemConfigCache.Cache.PartitionSync.SyncPolicy.PolicySwitch {
			PartitionSyncSingleTime()
		}
	}
}

func periodicRemoteBatchSync() {
}

func periodicRemoteSingleSync() {
}

func timingLocalBatchSync() {
	if config.SystemConfigCache.Cache.LocalBatchSync.SyncPolicy.TimingSync.Enable {
		if config.SystemConfigCache.Cache.LocalBatchSync.SyncPolicy.PolicySwitch {
			LocalBatchSyncSingleTime(true)
		}
	}

	if config.SystemConfigCache.Cache.LocalBatchSync.SyncPolicy.TimingSync.Enable {
		if config.SystemConfigCache.Cache.LocalBatchSync.SyncPolicy.PolicySwitch {
			nextTime := getNextTimeFromConfig(true, false, false)
			if nextTime == 0 {
				time.Sleep(61 * time.Second)
			}
			nextTime = getNextTimeFromConfig(true, false, false)
			notEnd := false
			StartPolicySync(nextTime, &notEnd, true, false, false, false)
		}
	}

}

func timingLocalSingleSync() {
	if config.SystemConfigCache.Cache.LocalSingleSync.SyncPolicy.TimingSync.Enable {
		LocalSingleSyncSingleTime(true)
	}

	if config.SystemConfigCache.Cache.LocalSingleSync.SyncPolicy.TimingSync.Enable {
		nextTime := getNextTimeFromConfig(false, false, false)
		notEnd := false
		StartPolicySync(nextTime, &notEnd, false, false, false, false)
	}
}

func timingPartitionSync() {
	if config.SystemConfigCache.Cache.PartitionSync.SyncPolicy.TimingSync.Enable {
		PartitionSyncSingleTime()
	}
	if config.SystemConfigCache.Cache.PartitionSync.SyncPolicy.TimingSync.Enable {
		nextTime := getNextTimeFromConfig(false, false, true)
		notEnd := false
		StartPolicySync(nextTime, &notEnd, false, false, false, false)
	}
}

func timingRemoteBatchSync() {
}

func timingRemoteSingleSync() {
}

func tickerWorker(ticker *time.Ticker, duration time.Duration, notEnd *bool, workerFunc func(), isPeriodic bool) {
	if duration == 0 {
		duration = time.Second
	}
	ticker = time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			if !*notEnd {
				go workerFunc()
				if !isPeriodic {
					ticker.Stop()
				}
			} else {
				return
			}
		}
	}
}

func StartPolicySync(duration time.Duration, notEnd *bool, isBatch bool, isRemote bool, isPartition bool, isPeriodic bool) {
	ticker, workerFunc := getTicker(isBatch, isRemote, isPeriodic, isPartition)
	if ticker != nil {
		ticker.Stop()
	}
	go tickerWorker(ticker, duration, notEnd, workerFunc, isPeriodic)
}

func PartitionSyncSingleTime() {
	if common.GetRunningFlag(common.TYPE_PARTITION) {
		return
	}
	sn, lock, startTime := preSyncSingleTime(common.TYPE_PARTITION, common.PARTITION_RUNNING)
	defer afterSyncSingleTime(common.TYPE_PARTITION, common.PARTITION_FORCE_DONE)
	GetTotalSize(&sn, config.SystemConfigCache.Cache.PartitionSync.SourcePath, true, lock)
	lock.Done()
	batchSyncFile(
		config.SystemConfigCache.Cache.PartitionSync.SourcePath,
		config.SystemConfigCache.Cache.PartitionSync.TargetPath,
		&sn,
		common.TYPE_PARTITION,
		false,
	)
	recordLog(common.TYPE_PARTITION,
		startTime,
		config.SystemConfigCache.Cache.PartitionSync.TargetPath,
		config.SystemConfigCache.Cache.PartitionSync.SourcePath,
	)
}

func getSingleTargetFile(sf *os.File, targetPath string) *os.File {
	tempTarget := ""
	tf, err := common.OpenFile(config.SystemConfigCache.Value().LocalSingleSync.TargetPath, true)

	if err != nil {
		if common.IsOpenDirError(err, config.SystemConfigCache.Value().LocalSingleSync.TargetPath) {
			sfInfo, _ := sf.Stat()
			tempTarget = config.SystemConfigCache.Value().LocalSingleSync.TargetPath + "/" + sfInfo.Name()
			tf, err = common.OpenFile(tempTarget, true)
			if err != nil {
				return tf
			}
		} else {
			return nil
		}
	}
	return tf
}

func localSyncSingleFile(msg interface{}, q *common.QWorker) {
	buf := make([]byte, 4096*2)
	common.SetCurrentSyncFile(q.SN, SYNC_RUNNING, q.PrivateFile.Name())
	for {
		n, err := q.PrivateFile.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			break
		}
		_, err = q.TargetFile.Write(buf[:n])
		if err != nil {
			sendNameToErrMap(q.SN, q.TargetFile.Name())
			break
		} else {
			addSizeToDoneMap(q.SN, uint64(n))
		}
	}
	syncFilePerm(q.PrivateFile, q.TargetFile)
}

func syncFilePerm(source *os.File, target *os.File) {
	sfInfo, _ := source.Stat()
	sfMode := sfInfo.Mode()
	err := target.Chmod(sfMode)
	if err != nil {
		log.Printf("Sync file perm err: %v", err)
	}
}

func closeAndCheckFile(w *common.QWorker) {
	if w.TargetFile != nil {
		tfName := w.TargetFile.Name()
		common.CloseFile(w.TargetFile)
		f, err := common.OpenFile(tfName, false)
		if err != nil {

		}
		w.TargetFile = f
	}
	if w.PrivateFile != nil {
		sfName := w.PrivateFile.Name()
		common.CloseFile(w.PrivateFile)
		f, err := common.OpenFile(sfName, false)
		if err != nil {

		}
		w.PrivateFile = f
	}
	common.SetCurrentSyncFile(w.SN, VERIFY_MD5, w.TargetFile.Name())
	if !CompareAndCacheMd5(w.PrivateFile, w.TargetFile, &w.SN, w.Md5CacheFlag) {
		AddBatchSyncError(w.PrivateFile.Name(), md5CheckError, w.SN)
	}
	if w.TargetFile != nil {
		common.CloseFile(w.TargetFile)
	}
	if w.PrivateFile != nil {
		common.CloseFile(w.PrivateFile)
	}
}

func GetLocalBatchProgress(sn string) float64 {
	resourceLock.Lock()
	defer resourceLock.Unlock()
	return float64(doneSizeMap[sn]) / float64(totalSizeMap[sn])
}

func GetBatchSyncError(sn string) []SyncFileError {
	errorMapLock.Lock()
	defer errorMapLock.Unlock()
	return syncErrorMap[sn]
}

func AddBatchSyncError(absPath string, reason string, sn string) {
	errorMapLock.Lock()
	defer errorMapLock.Unlock()
	node := SyncFileError{
		AbsPath: absPath,
		Reason:  reason,
	}
	syncErrorMap[sn] = append(syncErrorMap[sn], node)
}

func getTicker(isBatch bool, isRemote bool, isPeriodic bool, isPartition bool) (*time.Ticker, func()) {
	if isPeriodic {
		if isPartition {
			return periodicPartitionTicker, periodicPartitionSync
		} else if isRemote {
			if isBatch {
				return periodicRemoteBatchTicker, periodicRemoteBatchSync
			}
			return periodicRemoteSingleTicker, periodicRemoteSingleSync
		}
		if isBatch {
			return periodicLocalBatchTicker, periodicLocalBatchSync
		}
		return periodicLocalSingleTicker, periodicLocalSingleSync
	} else {
		if isPartition {
			return timingPartitionTicker, timingPartitionSync
		} else if isRemote {
			if isBatch {
				return timingRemoteBatchTicker, timingRemoteBatchSync
			}
			return timingRemoteSingleTicker, timingRemoteSingleSync
		}
		if isBatch {
			return timingLocalBatchTicker, timingLocalBatchSync
		}
		return timingLocalSingleTicker, timingLocalSingleSync
	}
}

func addSizeToDoneMap(sn string, size uint64) {
	resourceLock.Lock()
	defer resourceLock.Unlock()
	doneSizeMap[sn] = doneSizeMap[sn] + size
}

func getDoneSize(sn string) uint64 {
	resourceLock.Lock()
	defer resourceLock.Unlock()
	return doneSizeMap[sn]
}

func getTotalSize(sn string) uint64 {
	resourceLock.Lock()
	defer resourceLock.Unlock()
	return totalSizeMap[sn]
}

func addSizeToTotalMap(sn string, size uint64) {
	resourceLock.Lock()
	defer resourceLock.Unlock()
	totalSizeMap[sn] = totalSizeMap[sn] + size
}

func initSizeMap(sn string) {
	resourceLock.Lock()
	defer resourceLock.Unlock()
	zero := uint64(0)
	totalSizeMap[sn] = zero
	doneSizeMap[sn] = zero
}

func sendNameToErrMap(sn string, name string) {
	resourceLock.Lock()
	defer resourceLock.Unlock()
	batchErrMap[sn] = append(batchErrMap[sn], name)
}
