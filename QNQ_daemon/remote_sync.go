package QNQ

import (
	"context"
	"errors"
	"log/slog"
	"os"
)

// FileNode 只有根节点的name是路径
type FileNode struct {
	Name     string
	Children []*FileNode
}

type RemoteSync struct {
	UnimplementedRemoteSyncServer
}

func initRemoteSync() {
	for _, v := range ConfigCache.RemoteQNQ {
		err := addRPClient(v.Ip, v.Port)
		if err != nil {
			slog.Error("init remote QNQ target err", "err", err.Error())
		}
	}
}

func (rs *RemoteSync) AddTarget(ctx context.Context, req *RemoteTargetRequest) (*Result, error) {
	res := &Result{
		Code: ERR_CODE,
	}
	target := target{
		Ip:   req.Ip,
		Port: int(req.Port),
		Des:  "",
	}
	targetInCache := false
	for _, v := range ConfigCache.RemoteQNQ {
		if v.Ip == req.Ip {
			targetInCache = true
			break
		}
	}
	if _, ok := clientMap[req.Ip]; ok {
		loadErrToResult(res, errors.New("target QNQ is connected ,ip : "+req.Ip))
	} else if !targetInCache {
		err := addRPClient(target.Ip, target.Port)
		if err == nil {
			ConfigCache.RemoteQNQ = append(ConfigCache.RemoteQNQ, target)
			ConfigCache.notifyAll()
			res.Code = OK_CODE
		} else {
			loadErrToResult(res, err)
		}
	}
	return res, nil
}

func (rs *RemoteSync) DeleteTarget(ctx context.Context, req *RemoteTargetRequest) (*Result, error) {
	res := &Result{
		Code: ERR_CODE,
	}
	for i, v := range ConfigCache.RemoteQNQ {
		if v.Ip == req.Ip {
			ConfigCache.RemoteQNQ = append(ConfigCache.RemoteQNQ[:i], ConfigCache.RemoteQNQ[i+1:]...)
			ConfigCache.notifyAll()
			break
		}
	}
	err := deleteRPClient(req.Ip)
	if err != nil {
		loadErrToResult(res, err)
	}
	return res, nil
}

func (rs *RemoteSync) GetFileInfo(ctx context.Context, req *GetRemoteFileInfoRequest) (*FileInfoResult, error) {
	var err error
	result := resultPool.Get().(*Result)
	res := &FileInfoResult{
		Result:    result,
		FileInfos: make([]*FileInfoResult_FileInfo, 0),
	}
	root := NewFileNode(req.Path)
	GetFileTree(root, "", 0, 1)
	res.Len = int64(len(root.Children))
	res.FileInfos = make([]*FileInfoResult_FileInfo, 0, res.Len)
	for _, child := range root.Children {
		fi := conversionNode2FileInfo(child, req.Path)
		res.FileInfos = append(res.FileInfos, fi)
	}
	res.Result.Code = OK_CODE

	return res, err
}

func (rs *RemoteSync) InputSync(ctx context.Context, req *RemoteSyncRequest) (*Result, error) {
	return nil, nil
}

func (rs *RemoteSync) OutputSync(ctx context.Context, req *RemoteSyncRequest) (*Result, error) {
	return nil, nil
}

func conversionNode2FileInfo(node *FileNode, rootPath string) *FileInfoResult_FileInfo {
	fullPath := rootPath + sparator + node.Name
	res := &FileInfoResult_FileInfo{
		Name: node.Name,
	}
	fInfo, err := os.Stat(fullPath)
	if err == nil {
		res.Prem = uint32(fInfo.Mode())
		res.ModifyDate = fInfo.ModTime().Format(timeFormat)
		res.IsDir = fInfo.IsDir()
		if !res.IsDir {
			res.Size = uint64(fInfo.Size())
		}
	}
	return res
}

func NewFileNode(name string) *FileNode {
	return &FileNode{
		Name:     name,
		Children: make([]*FileNode, 0),
	}
}

func findFileTreeNode(root *FileNode, targetPath string, currentPath string) *FileNode {
	if root == nil {
		return nil
	}
	var newPath string
	if "" != currentPath {
		newPath = currentPath + sparator + root.Name
	} else {
		newPath = root.Name
	}

	if newPath == targetPath[:len(newPath)] {
		if newPath == targetPath {
			return root
		}
		var res *FileNode
		for _, child := range root.Children {
			res = findFileTreeNode(child, targetPath, newPath)
			if res != nil {
				return res
			}
		}
	}
	return nil
}

// GetFileTree root节点需要提前初始化，maxDeep 最大层深，-1为不限制层深, currentPath传空字符串
func GetFileTree(root *FileNode, lastPath string, deep, maxDeep int) {
	if root == nil || (maxDeep != -1 && deep > maxDeep) {
		return
	}
	var currentPath string
	if "" != lastPath {
		currentPath = lastPath + sparator + root.Name
	} else {
		currentPath = root.Name
	}
	fInfo, err := os.Stat(currentPath)
	if err != nil {
		slog.Error(err.Error() + ": " + currentPath)
		return
	}
	if fInfo.IsDir() {
		f, err := openDir(currentPath)
		if err != nil {
			return
		}
		children, err := f.Readdir(-1)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		for _, child := range children {
			node := NewFileNode(child.Name())
			root.Children = append(root.Children, node)
			GetFileTree(node, currentPath, deep+1, maxDeep)
		}
	}

}

func GetRemoteDir(ip string) {

}
