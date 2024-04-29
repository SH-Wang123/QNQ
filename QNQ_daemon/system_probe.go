package QNQ

import (
	"context"
	"github.com/shirou/gopsutil/disk"
	"log/slog"
	"runtime"
)

var ipBitMap uint32

const localhost = "127.0.0.1"

type partitionStat struct {
	device      string
	mountPoint  string
	opts        string
	fsType      string
	total       uint64
	free        uint64
	usedPercent float64
}

type sysInfo struct {
	hostOs         string
	hostArch       string
	partitionStats []*partitionStat
}

func initSysInfo() {
	ConfigCache.sysInfo = sysInfo{
		hostOs:   runtime.GOOS,
		hostArch: runtime.GOARCH,
	}
	loadDiskInfo()
}

func loadDiskInfo() {
	parts, err := disk.Partitions(true)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	ConfigCache.sysInfo.partitionStats = make([]*partitionStat, 0)
	for _, part := range parts {
		p := &partitionStat{}
		p.fsType = part.Fstype
		p.device = part.Device
		p.mountPoint = part.Mountpoint
		p.opts = part.Opts
		usageStat, err := disk.Usage(part.Device)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		p.total = usageStat.Total
		p.free = usageStat.Free
		p.usedPercent = usageStat.UsedPercent
		ConfigCache.sysInfo.partitionStats = append(ConfigCache.sysInfo.partitionStats, p)
	}
}

type SysProbe struct {
	UnimplementedSysProbeServer
}

func (s *SysProbe) GetSysInfo(ctx context.Context, req *GetSysInfoRequest) (*SysInfoResult, error) {
	var err error
	var res *SysInfoResult
	result := resultPool.Get().(*Result)
	res = &SysInfoResult{
		Result: result,
	}
	res.HostArch = ConfigCache.sysInfo.hostArch
	res.HostOs = ConfigCache.sysInfo.hostOs
	res.QnqVersion = ConfigCache.Version
	for _, p := range ConfigCache.sysInfo.partitionStats {
		resPart := &SysInfoResultPartitionStat{
			Device:      p.device,
			MountPoint:  p.mountPoint,
			Opts:        p.opts,
			FsType:      p.fsType,
			Total:       p.total,
			Free:        p.free,
			UsedPercent: p.usedPercent,
		}
		res.PartitionStats = append(res.PartitionStats, resPart)
	}
	res.Result.Code = OK_CODE
	slog.Info("GetSysInfo", "res", res)
	return res, err
}

func (s *SysProbe) GetPerfProbe(ctx context.Context, req *GetPerfProbeRequest) (*PerfProbeResult, error) {
	result := resultPool.Get().(*Result)
	res := &PerfProbeResult{
		Result: result,
	}

	return res, nil
}
