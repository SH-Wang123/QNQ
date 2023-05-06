package rest

import (
	"github.com/gin-gonic/gin"
	"window_handler/common"
	"window_handler/worker"
)

func getFileRootMap(c *gin.Context) {
	absPath := c.Query("absPath")
	pointPath := c.Query("pointPath")
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
