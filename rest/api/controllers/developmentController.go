package controllers

import (
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/daoPostgres/daoQuery/enums/daoOperator"
	"Quote/db/daoPostgres/daos/apertureDao"
	"Quote/db/daoPostgres/daos/fileDao"
	"Quote/db/daoPostgres/daos/profileDao"
	"Quote/db/daoPostgres/daos/quoteDao"
	"Quote/db/entity/entities/apertureEntity"
	"Quote/db/entity/entities/fileEntity"
	"Quote/db/entity/entities/quoteEntity"
	"Quote/model/enums/fileTypes"
	"Quote/pdf/supplierPdf"
	"Quote/utils/imageUtils"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	ctx := context.Background()
	quoteId := uint64(2)

	selector := daoQuery.Selector{
		SubItems: []string{"clients"},
	}
	quoteI, _, err := quoteDao.GetInstance().ReadByID(ctx, quoteId, &selector)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	quote := quoteI.(*quoteEntity.QuoteEntity)

	var mainApertureList []*apertureEntity.ApertureEntity
	// Get the list of apertures
	{
		noImageBytes, err := imageUtils.FromFilePath("./images/no-image.jpg")
		if err != nil {
			fmt.Printf("%s", err.Error())
			return
		}
		noImageEntity := &fileEntity.FileEntity{
			Data: noImageBytes,
			Type: int(fileTypes.Jpg),
		}

		p := daoQuery.NewPagination()
		p.Applies = false

		selector := daoQuery.Selector{
			SubItems: []string{"colors", "glasses", "profiles"},
		}
		err = apertureDao.GetInstance().ReadMainByQuoteId(ctx, quoteId, p, &selector)
		if err != nil {
			fmt.Printf("%s", err.Error())
			return
		}

		for _, apertureI := range p.Data.([]any) {
			aperture := apertureI.(*apertureEntity.ApertureEntity)

			// Fill the aperture's image
			{
				aperture.Image = noImageEntity
				if aperture.ImageId != nil {
					imgEntityI, found, err := fileDao.GetInstance().ReadByID(ctx, *aperture.ImageId, nil)
					if err != nil {
						fmt.Printf("%s", err.Error())
						return
					}

					if found {
						imgEntity := imgEntityI.(*fileEntity.FileEntity)
						aperture.Image = imgEntity
					}
				}
			}

			mainApertureList = append(mainApertureList, aperture)
		}

		// apertureMap := map[uint64]*apertureEntity.ApertureEntity{}
		// orderedKeys := []uint64{}
		// for _, apertureI := range p.Data.([]any) {
		// 	aperture := apertureI.(*apertureEntity.ApertureEntity)

		// 	if aperture.ParentApertureID == nil {
		// 		apertureMap[aperture.ID] = aperture
		// 		orderedKeys = append(orderedKeys, aperture.ID)
		// 	} else {
		// 		parentAperture := apertureMap[*aperture.ParentApertureID]
		// 		parentAperture.ChildApertures = append(parentAperture.ChildApertures, aperture)
		// 	}
		// }

		// for _, key := range orderedKeys {
		// 	mainApertureList = append(mainApertureList, apertureMap[key])
		// }
	}

	var datasheetList [][]byte
	// Get the list of datasheets
	{
		apertureDao := apertureDao.GetInstance()
		profileDao := profileDao.GetInstance()
		fileDao := fileDao.GetInstance()
		tableName := apertureDao.TableName()
		profileTableName := profileDao.TableName()
		fileTableName := fileDao.TableName()

		tx := apertureDao.GetBasicTransaction(ctx, apertureDao.NewEntity, false)
		tx = tx.Select("profiles.datasheet_id as id")

		tx.Joins(fmt.Sprintf("left join %s on %s.profile_id = %s.id", profileTableName, tableName, profileTableName))
		tx.Joins(fmt.Sprintf("left join %s on %s.datasheet_id = %s.id", fileTableName, profileTableName, fileTableName))
		tx.Where("files.type = ?", fileTypes.Pdf)
		tx.Scopes(
			apertureDao.ByQuoteId(daoOperator.Equals, quoteId),
			apertureDao.ByProfileId(daoOperator.IsNot, nil),
			apertureDao.ByDatasheetId(daoOperator.IsNot, nil),
		)
		tx.Group("profiles.id, files.type")

		itemListMap := []map[string]any{}
		if err := tx.Find(&itemListMap).Error; err != nil {
			fmt.Printf("%s", err.Error())
			return
		}

		for _, item := range itemListMap {
			id := item["id"].(uint64)
			fileI, found, err := fileDao.ReadByID(ctx, id, nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return
			}

			if found {
				file := fileI.(*fileEntity.FileEntity)
				datasheetList = append(datasheetList, file.Data)
			}
		}
	}

	pdfBytes, err := supplierPdf.Generate(quote, mainApertureList, datasheetList, nil)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	err = os.WriteFile("test.pdf", pdfBytes, 0644) // Creates the PDF file
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
}
