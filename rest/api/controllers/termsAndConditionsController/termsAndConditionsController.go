package termsAndConditionsController

import (
	"Quote/db/daoPostgres/daos/fileDao"
	"Quote/db/daoPostgres/daos/termsAndConditionsDao"
	"Quote/logger"
	"Quote/rest/api/apiUtils"
	"Quote/rest/api/baseController"
	"Quote/service/fileService"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const withAuthorisation = true

const fileName = "terms-and-conditions.pdf"

// Exists reports whether a terms and conditions PDF is stored in the database.
func Exists(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	_, found, err := termsAndConditionsDao.GetInstance().Get(ctx)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to check the terms and conditions: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": found})
}

// Download returns the terms and conditions PDF so that it can be viewed in the browser.
func Download(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	terms, found, err := termsAndConditionsDao.GetInstance().Get(ctx)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch the terms and conditions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "terms and conditions not available", http.StatusNotFound)
		return
	}

	data, found, err := fileService.GetSingleFile(ctx, terms.FileID)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch the terms and conditions file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "terms and conditions not available", http.StatusNotFound)
		return
	}

	c.Header("Content-Disposition", "inline; filename="+fileName)
	c.Data(http.StatusOK, "application/pdf", data)
}

// Upload stores the submitted PDF (multipart field "file") as the terms and conditions, replacing any existing one.
// It responds with the id of the stored file.
func Upload(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	authUserId, err := baseController.GetAuthenticatedId(c)
	if withAuthorisation && err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnauthorized)
		return
	}

	files, err := apiUtils.HttpMultipartForm(c, "file")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}
	if len(files) != 1 {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "exactly one file must be submitted in the 'file' field", http.StatusBadRequest)
		return
	}
	file := files[0]

	// Only PDFs are accepted: check the declared type and the file signature
	if file.Header.Get("Content-Type") != "application/pdf" {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "only PDF files (application/pdf) can be uploaded", http.StatusUnsupportedMediaType)
		return
	}
	{
		opened, err := file.Open()
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to open the file: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		signature := make([]byte, 5)
		_, err = io.ReadFull(opened, signature)
		opened.Close()
		if err != nil || string(signature) != "%PDF-" {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), "the submitted file is not a valid PDF", http.StatusUnsupportedMediaType)
			return
		}
	}

	fileId, err := fileService.AddSingleFile(ctx, authUserId, file)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}

	replacedFileId, err := termsAndConditionsDao.GetInstance().Set(ctx, fileId, authUserId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to save the terms and conditions: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// The replaced PDF is no longer referenced, so soft delete it
	if replacedFileId != 0 {
		if _, err := fileDao.GetInstance().SoftDelete(ctx, replacedFileId, authUserId); err != nil {
			logger.GetInstance().Warnf("failed to soft delete the replaced terms and conditions file %d: %s", replacedFileId, err.Error())
		}
	}

	c.JSON(http.StatusOK, fileId)
}
