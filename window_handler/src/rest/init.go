package rest

import (
	"github.com/gin-gonic/gin"
	"log"
	"window_handler/common"
)

func init() {
	go func() {
		r := gin.Default()

		r.GET("/fileRootMap", getFileRootMap)

		r.GET("/disk/info", getDiskInfo)

		r.GET("/debug/connect", getConnect)

		err := r.Run(common.QNQ_TARGET_REST_PORT)
		if err != nil {
			log.Printf("create router err : %v", err)
			return
		}
	}()
}
