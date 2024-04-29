//go:build windows

package QNQ

import (
	"io"
	"log/slog"
	"os"
)

var sparator = "\\"

func syncFile(source string, target string, t *Task) *TaskResult {
	res := &TaskResult{
		Code: OK_CODE,
	}
	sf, errSf := os.Open(source)
	defer closeFile(sf)
	srcInfo, _ := sf.Stat()
	if errSf != nil {
		res.SetCode(SourceFileOpenFailed)
		return res
	}
	preCheck := localSyncPreCheck(srcInfo.IsDir(), target)
	if preCheck != OK_CODE {
		res.SetCode(preCheck)
		return res
	}
	if srcInfo.IsDir() {
		children, _ := sf.Readdir(-1)
		for _, child := range children {
			currentSource := source + "\\" + child.Name()
			currentTarget := target + "\\" + child.Name()
			if child.IsDir() {
				cRes := syncFile(currentSource, currentTarget, t)
				if cRes.Code != OK_CODE {
					return cRes
				}
			} else {
				currentSf, errSf := OpenFile(currentSource)
				if errSf != nil {
					res.SetCode(SourceFileOpenFailed)
					return res
				}
				currentTf, errTf := OpenFile(currentTarget)
				if errTf != nil {
					res.SetCode(TargetFileOpenFailed)
					return res
				}
				if res := syncSingleFile(currentSf, currentTf, t); !res {
					t.Status = ERR_CODE
				}
			}

		}
	} else {
		tf, errTF := os.Open(target)
		if errTF != nil {
			res.SetCode(TargetFileOpenFailed)
			return res
		} else {
			tarInfo, _ := tf.Stat()
			if tarInfo.IsDir() {
				tf, errTF = OpenFile(target + "\\" + srcInfo.Name())
				if errTF != nil {
					res.SetCode(TargetFileOpenFailed)
					return res
				}
			}
		}
		if res := syncSingleFile(sf, tf, t); !res {
			t.Status = ERR_CODE
		}
	}

	return res
}

func syncSingleFile(sf *os.File, tf *os.File, t *Task) bool {
	buf := make([]byte, 4096)
	for {
		n, err := sf.Read(buf)
		if err != nil && err != io.EOF {
			return true
		}
		if n == 0 {
			break
		}
		_, err = tf.Write(buf[:n])
		t.Probe.doneChan <- uint64(n)
		if err != nil {
			break
		}
	}

	defer func(tf *os.File) {
		go func() {
			err := tf.Sync()
			if err != nil {
				slog.Error(err.Error())
				return
			}
			err = tf.Close()
			if err != nil {
				slog.Error(err.Error())
				return
			}
		}()
	}(tf)

	return false
}

func getTotalSize(source string, t *Task) {
	f, err := os.Stat(source)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if f.IsDir() {
		sf, err := openDir(source)
		if err != nil {

		}
		children, _ := sf.Readdir(-1)
		for _, child := range children {
			getTotalSize(source+sparator+child.Name(), t)
		}
	} else {
		t.Probe.totalChan <- uint64(f.Size())
	}
}
