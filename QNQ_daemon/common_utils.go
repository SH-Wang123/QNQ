package QNQ

import (
	"errors"
	"log"
	"log/slog"
	"os"
	"reflect"
	"sync"
	"syscall"
	"time"
)

type capacityUnit uint64

var resultPool = sync.Pool{
	New: func() any {
		return &Result{
			Code:    ERR_CODE,
			Message: "",
			TaskId:  0,
		}
	},
}

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

type ProgressProbe struct {
	totalSize uint64
	doneSize  uint64
	totalChan chan uint64
	doneChan  chan uint64
	overChan  chan struct{}
}

func (p *ProgressProbe) watch() {
	go func() {
		for {
			select {
			case v, ok := <-p.totalChan:
				if ok {
					p.totalSize += v
				}
			case v, ok := <-p.doneChan:
				if ok {
					p.doneSize += v
				}
			case _, ok := <-p.overChan:
				if ok {
					return
				}
			}
		}
	}()
}

func (p *ProgressProbe) GetProgress() float32 {
	if float32(p.totalSize) == 0 {
		return 0
	}
	return float32(p.doneSize) / float32(p.totalSize)
}

func (p *ProgressProbe) GetRate() float32 {
	startSize := p.doneSize
	time.Sleep(1 * time.Second)
	return (float32(p.doneSize) - float32(startSize)) / MB
}

func NewProgressProbe() *ProgressProbe {
	return &ProgressProbe{
		totalChan: make(chan uint64),
		doneChan:  make(chan uint64),
		overChan:  make(chan struct{}),
	}
}

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteFileOrDir(path string) {
	exist, err := IsExist(path)
	f, _ := OpenFile(path)
	fChild, _ := f.Readdir(-1)
	closeFile(f)
	if err != nil {
		log.Printf("get file error : %v", err)
	}
	if exist {
		if len(fChild) == 0 {
			err = os.Remove(path)
		} else {
			err = os.RemoveAll(path)
		}
	}
	if err != nil {
		slog.Error("delete file error : %v", err)
	}
}

func OpenFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return f, nil
}

func closeFile(fs ...*os.File) {
	for _, f := range fs {
		if f == nil {
			continue
		}
		err := f.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func openDir(path string) (*os.File, error) {
	f, err := os.Open(path)
	if errors.Is(err, syscall.ERROR_FILE_NOT_FOUND) {
		err = os.Mkdir(path, 0777)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		} else {
			f, err = os.Open(path)
			if err != nil {
				slog.Error(err.Error())
				return nil, err
			}
		}
	}
	return f, err
}

func localSyncPreCheck(sourceIsDir bool, target string) int {
	var errTF error
	var tf *os.File
	if sourceIsDir {
		tf, errTF = openDir(target)
	} else {
		tf, errTF = os.Open(target)
	}
	if errTF != nil {
		defer closeFile(tf)
		return TargetFileOpenFailed
	} else {
		tarInfo, _ := tf.Stat()
		if !tarInfo.IsDir() && sourceIsDir {
			return DirSync2FileError
		}
	}

	return OK_CODE
}

func isLocalTarget(ip string) bool {
	return localhost == ip
}

// hasField Check if there is a fieldName in the struct and set the targetValue
func hasField(obj any, fieldName string, targetValue any) (bool, any) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return false, nil
	}
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if fieldName == field.Name {
			value := val.Field(i).Interface()
			if targetValue != nil {
				val.Field(i).Set(reflect.ValueOf(targetValue))
			}
			return true, value
		}
	}
	return false, nil
}

func reflectMethod(entity any, name string, params ...any) ([]any, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("reflectMethod panic", "recover", r)
		}
	}()
	var res []any
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, errors.New("entity is not a struct")
	}
	methodValue := reflect.ValueOf(entity).MethodByName(name)
	if methodValue.IsValid() {
		paramsValue := make([]reflect.Value, len(params))
		for i := 0; i < len(params); i++ {
			paramsValue[i] = reflect.ValueOf(params[i])
		}
		resValue := methodValue.Call(paramsValue)
		res = make([]any, len(resValue))
		for i := 0; i < len(resValue); i++ {
			res[i] = resValue[i].Interface()
		}
	} else {
		return nil, errors.New("method not declaration")
	}
	return res, nil
}
