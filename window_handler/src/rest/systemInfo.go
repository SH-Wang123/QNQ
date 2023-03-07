package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"window_handler/common"
	"window_handler/worker"
)

func getFileRootMap(c *gin.Context) {
	abstractPath := c.Query("abstractPath")
	anchorPointPath := c.Query("anchorPointPath")
	startNode := &worker.FileNode{
		IsDirectory:     true,
		HasChildren:     true,
		AbstractPath:    fmt.Sprintf("%v", abstractPath),
		AnchorPointPath: fmt.Sprintf("%v", anchorPointPath),
		HeadFileNode:    nil,
	}
	worker.GetFileTree(startNode, false)
	dataM := make(map[string][]string)
	worker.GetFileTreeMap(startNode, &dataM)
	rsp := common.NewQResponse(RESPONSE_CODE_OK, dataM)
	c.JSON(HTTP_OK, rsp)
}

func getDiskInfo(c *gin.Context) {
	rsp := common.NewQResponse(RESPONSE_CODE_OK, worker.DiskPartitionsCache)
	if worker.DiskPartitionsCache != nil {
		c.JSON(HTTP_OK, rsp)
	} else {
		c.JSON(HTTP_ERROR, rsp)
	}
}

func getConnect(c *gin.Context) {
	rsp := common.NewQResponse(RESPONSE_CODE_OK, "ok")
	c.JSON(HTTP_OK, rsp)
}
