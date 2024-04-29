package QNQ

import (
	"context"
)

type LocalSync struct {
	UnimplementedLocalSyncServer
}

func (s *LocalSync) SingleSync(ctx context.Context, req *LocalSyncRequest) (*Result, error) {
	res := resultPool.Get().(*Result)
	task := NewLocalSyncTask(req.GetSource(), req.GetTarget())
	res.TaskId = task.Id
	ConfigCache.setLocalSingleSync(req.GetSource(), req.GetTarget())
	CommitTask(task)
	return res, nil
}
func (s *LocalSync) BatchSync(ctx context.Context, req *LocalSyncRequest) (*Result, error) {
	res := &Result{}
	task := NewLocalSyncTask(req.GetSource(), req.GetTarget())
	res.TaskId = task.Id
	ConfigCache.setLocalBatchSync(req.GetSource(), req.GetTarget())
	CommitTask(task)
	return res, nil
}

func (s *LocalSync) PartitionSync(ctx context.Context, req *LocalSyncRequest) (*Result, error) {
	res := &Result{}
	task := NewLocalSyncTask(req.GetSource(), req.GetTarget())
	res.TaskId = task.Id
	CommitTask(task)
	return res, nil
}

func NewLocalSyncTask(sourcePath string, targetPath string) *Task {
	params := make(map[string]string)
	params["sourcePath"] = sourcePath
	params["targetPath"] = targetPath
	task := NewTask(localSyncExec, localSyncCancel, "params", params)
	task.Context.Value(params)
	return task
}

func localSyncExec(t *Task) {
	params := t.Context.Value("params").(map[string]string)
	getTotalSize(params["sourcePath"], t)
	t.Result = syncFile(params["sourcePath"], params["targetPath"], t)
	t.Result.Message = GetMsg(t.Result.Code)
}

func localSyncCancel(t *Task) {

}
