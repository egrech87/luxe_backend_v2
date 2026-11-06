package fileController

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daos/fileDao"
	"Quote/db/entity/entities/fileEntity"
	"Quote/logger"
	"Quote/model/enums/fileTypes"
	"Quote/model/models/fileModel"
	"Quote/rest/api/apiUtils"
	"Quote/rest/api/baseController"
	"Quote/service/fileService"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const withAuthorisation = true

func GetDao() baseDao.BaseDaoI {
	return fileDao.GetInstance()
}

func AddFile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	authUserId, err := baseController.GetAuthenticatedId(c)
	if withAuthorisation && err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnauthorized)
		return
	}

	file, err := apiUtils.HttpMultipartForm(c, "file")
	{
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
			return
		}
		if len(file) == 0 {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), "File not submitted to be uploaded", http.StatusNotFound)
			return
		}
		if len(file) > 1 {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), "More than 1 file detected to be uploaded. Only one file can be uploaded", http.StatusNotFound)
			return
		}
	}

	fileId, err := fileService.AddSingleFile(ctx, authUserId, file[0])
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, fileId)
}

func GetById(c *gin.Context) {
	baseController.GetById(c, "file_id", GetDao())
}

func DownloadFile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	fileName := "myfile.txt"

	id, err := apiUtils.HttpPathUInt64(c, "file_id")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	itemI, found, err := fileDao.GetInstance().ReadByID(ctx, *id, nil)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch item by id: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "item not available", http.StatusNotFound)
		return
	}

	fileEntity := itemI.(*fileEntity.FileEntity)
	fileModel := fileModel.New()
	fileModel.SetEntity(fileEntity)

	contentType := ""
	switch *fileModel.GetType() {
	case fileTypes.Png:
		contentType = "image/png"
	case fileTypes.Jpg:
		contentType = "image/jpeg"
	case fileTypes.Gif:
		contentType = "image/gif"
	case fileTypes.Pdf:
		contentType = "application/pdf"
	case fileTypes.WordDoc:
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case fileTypes.Excel:
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", contentType)
	c.Header("Accept-Length", fmt.Sprintf("%d", len(fileModel.GetData())))
	c.Writer.Write(fileModel.GetData())
}
