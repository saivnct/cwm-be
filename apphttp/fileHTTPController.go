package apphttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"sol.go/cwm/dao"
	"sol.go/cwm/s3Handler"
)

type FileHTTPController struct {
}

func (sv *FileHTTPController) getFile(ctx *gin.Context) {
	msgId := ctx.Param("msgId")
	fileId := ctx.Param("fileId")

	s3FileInfo, err := dao.GetS3FileInfoDAO().FindByFileId(ctx, fileId)
	if s3FileInfo == nil {
		log.Println("Not found s3FileInfo", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "not found",
		})
		return
	}

	if s3FileInfo.MsgId != msgId {
		log.Println("Invalid msgId")

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "invalid info",
		})
		return
	}

	s3Object, err := s3Handler.GetS3FileStore().GetObject(s3FileInfo.FileName)
	if err != nil {
		log.Println("Get s3 Object err", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "not found object",
		})

		return
	}

	if s3Object.Body != nil {
		defer s3Object.Body.Close()
	}
	// If there is no content length, it is a directory
	if s3Object.ContentLength == nil {
		return
	}

	//log.Println("Content-Type", *s3Object.ContentType)

	ctx.Header("Content-Type", *s3Object.ContentType)
	ctx.Header("Content-Length", fmt.Sprintf("%d", *s3Object.ContentLength))
	ctx.Status(http.StatusOK)
	io.Copy(ctx.Writer, s3Object.Body)
}
