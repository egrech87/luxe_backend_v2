package fileService

import (
	"Quote/db/daoPostgres/daos/fileDao"
	"Quote/db/entity/entities/fileEntity"
	"Quote/model/enums/fileTypes"
	"Quote/model/models/fileModel"
	"Quote/utils/imageUtils"
	"context"
	"mime/multipart"

	"github.com/pkg/errors"
)

func AddSingleFile(ctx context.Context, authUserId uint64, file *multipart.FileHeader) (fileId uint64, err error) {
	contentType := file.Header.Get("Content-Type")

	fileBytesList, err := imageUtils.FromMultipartFileHeader([]*multipart.FileHeader{file})
	if err != nil {
		return 0, err
	}

	fileModel := fileModel.New()
	fileModel.SetEntity(fileEntity.New())
	fileModel.SetData(fileBytesList[0])
	switch contentType {
	case "image/png":
		fileModel.SetType(fileTypes.Png)
	case "image/jpeg":
		fileModel.SetType(fileTypes.Jpg)
	case "image/gif":
		fileModel.SetType(fileTypes.Gif)
	case "application/pdf":
		fileModel.SetType(fileTypes.Pdf)
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		fileModel.SetType(fileTypes.WordDoc)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		fileModel.SetType(fileTypes.Excel)
	}

	if err = fileModel.ValidateType(ctx); err != nil {
		return 0, errors.Errorf("%s not supported", contentType)
	}

	if err = fileModel.Validate(ctx); err != nil {
		return 0, err
	}

	err = fileDao.GetInstance().Create(ctx, fileModel.GetEntity(), authUserId)
	if err != nil {
		return 0, err
	}

	return fileModel.GetID(), nil
}

func GetSingleFile(ctx context.Context, fileId uint64) (file []byte, found bool, err error) {
	fileEntityI, found, err := fileDao.GetInstance().ReadByID(ctx, fileId, nil)
	if err != nil || !found {
		return []byte{}, found, err
	}

	fileEntity := fileEntityI.(*fileEntity.FileEntity)

	return fileEntity.GetData(), found, err
}
