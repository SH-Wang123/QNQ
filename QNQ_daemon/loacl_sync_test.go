package QNQ

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"
)

type ret struct {
	age int
}

func (r *ret) TestFunc(i int, nr ret) (*ret, error) {
	return &nr, errors.New("test error")
}

func outputError(t *testing.T, task *Task) {
	t.Error("code : " + fmt.Sprint(task.Result.Code) + "  message : " + task.Result.Message)
}

func TestLocalBatchSync(t *testing.T) {
	sourceInput := "E:\\apache-tomcat-8.5.84-windows-x64"
	targetInput := "F:\\qnq_t"
	task := NewLocalSyncTask(sourceInput, targetInput)
	CommitTask(task)
	pro := float32(0)
	for pro < 1 {
		pro = task.Probe.GetProgress()
		t.Log(pro)
		t.Log(task.Probe.GetRate())
	}
	WaitTaskSchedulerStop()
}

func TestLocalBatchSyncTask(t *testing.T) {
	for i := 0; i < 10; i++ {
		sourceInput := "E:\\apache-tomcat-8.5.84-windows-x64"
		targetInput := "F:\\qnq_t"
		task := NewLocalSyncTask(sourceInput, targetInput)
		CommitTask(task)
	}
	WaitTaskSchedulerStop()
}

func TestLocalSingleSync(t *testing.T) {
	sourceInput := "E:\\httpd-2.4.54.tar.bz2"
	targetInput := "F:\\qnq_t"
	task := NewLocalSyncTask(sourceInput, targetInput)
	task.Execute()
}

func TestLocalSingleSyncTask(t *testing.T) {
	sourceInput := "E:\\apache-tomcat-8.5.84-windows-x64"
	targetInput := "F:\\qnq_t"
	task0 := NewLocalSyncTask(sourceInput, targetInput)
	CommitTask(task0)

	sourceInput = "F:\\qnq_t"
	targetInput = "E:\\httpd-2.4.54.tar.bz2"
	task1 := NewLocalSyncTask(sourceInput, targetInput)
	CommitTask(task1)
	WaitTaskSchedulerStop()
	if task0.Result.Code != OK_CODE {
		outputError(t, task0)
	}
	if task1.Result.Code != DirSync2FileError {
		outputError(t, task1)
	}

}

func TestFileTree(t *testing.T) {
	root := NewFileNode("E:\\apache-tomcat-8.5.84-windows-x64")
	GetFileTree(root, "", 0, -1)
	nroot := findFileTreeNode(root, "E:\\apache-tomcat-8.5.84-windows-x64\\apache-tomcat-8.5.84\\bin", "")
	if nroot.Name != "bin" {
		t.Error("find error node : " + nroot.Name)
	}
}

func TestSlog(t *testing.T) {
	slog.Info("a", "b", "c", "d", "e")
	slog.Error("a", "b", "c")
}

func TestReflectMethod(t *testing.T) {
	ru := &ret{
		age: 10,
	}
	res, err := reflectMethod(ru, "TestFunc", 1, *ru)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(fmt.Sprint(res))
}
