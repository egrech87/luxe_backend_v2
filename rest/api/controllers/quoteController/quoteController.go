package quoteController

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/daoPostgres/daoQuery/enums/daoOperator"
	"Quote/db/daoPostgres/daos/apertureDao"
	"Quote/db/daoPostgres/daos/fileDao"
	"Quote/db/daoPostgres/daos/glassDao"
	"Quote/db/daoPostgres/daos/profileDao"
	"Quote/db/daoPostgres/daos/quoteDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/apertureEntity"
	"Quote/db/entity/entities/fileEntity"
	"Quote/db/entity/entities/quoteEntity"
	"Quote/logger"
	"Quote/model/baseModel"
	"Quote/model/enums/fileTypes"
	"Quote/model/models/quoteModel"
	"Quote/pdf/clientPdf"
	"Quote/pdf/supplierPdf"
	"Quote/rest/api/apiUtils"
	"Quote/rest/api/baseController"
	"Quote/service/termsAndConditionsService"
	"Quote/utils/imageUtils"
	"Quote/utils/pathUtils"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const withAuthorisation = true

func GetNewEntity() baseEntity.BaseEntityI {
	return quoteEntity.New()
}
func GetNewModel() baseModel.BaseModelI {
	return quoteModel.New()
}
func GetDao() baseDao.BaseDaoI {
	return quoteDao.GetInstance()
}
func GetSystemFields() []string {
	return append(baseController.GetSystemFields(), "total_area", "total_cost", "total_apertures")
}

func Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	authUserId, err := baseController.GetAuthenticatedId(c)
	if withAuthorisation && err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnauthorized)
		return
	}

	model := quoteModel.New()
	dao := quoteDao.GetInstance()

	// Get the model from the request
	{
		inputJson, err := apiUtils.HttpBodyJson(c, true)
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
			return
		}
		baseController.CleanInput(inputJson, GetSystemFields())
		model.FillFromJson(inputJson)

		isNumberCustom, err := apiUtils.HttpQueryBool(c, "is_number_custom", false)
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
			return
		}

		if isNumberCustom != nil && !*isNumberCustom {
			model.Number, err = dao.ReadNextNumber(ctx)
			if err != nil {
				apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := model.Validate(ctx); err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusConflict)
		return
	}

	err = dao.Create(ctx, model.GetEntity(), authUserId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.GetEntity().GetID())
}

func Update(c *gin.Context) {
	baseController.Update(c, "quote_id", GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func SoftDelete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	authUserId, err := baseController.GetAuthenticatedId(c)
	if withAuthorisation && err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnauthorized)
		return
	}

	id, err := apiUtils.HttpPathUInt64(c, "quote_id")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	deletedAt, err := quoteDao.GetInstance().SoftDelete(ctx, *id, authUserId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to delete the quote: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = apertureDao.GetInstance().SoftDeleteByQuoteId(ctx, *id, deletedAt, authUserId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to delete the apertures related to the quote: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func GetById(c *gin.Context) {
	baseController.GetById(c, "quote_id", GetDao())
}

func Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	dao := GetDao()

	selector, err := apiUtils.HttpHeaderSelector(c)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	pagination, err := apiUtils.HttpQueryPagination(c)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}
	if pagination.Meta.Sort == "" {
		pagination.Meta.Sort = fmt.Sprintf("%s.created_at DESC", dao.TableName())
	}

	err = dao.Read(ctx, pagination, selector)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch items: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	itemsMap := []map[string]any{}
	for _, item := range pagination.Data.([]any) {
		currentItem := item.(*quoteEntity.QuoteEntity)
		itemsMap = append(itemsMap, currentItem.ToJson())
	}
	pagination.Data = itemsMap

	c.JSON(http.StatusOK, pagination)
}

func GetDatasheets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	id, err := apiUtils.HttpPathUInt64(c, "quote_id")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	apertureDao := apertureDao.GetInstance()
	profileDao := profileDao.GetInstance()
	fileDao := fileDao.GetInstance()
	tableName := apertureDao.TableName()
	profileTableName := profileDao.TableName()
	fileTableName := fileDao.TableName()

	tx := apertureDao.GetBasicTransaction(ctx, apertureDao.NewEntity, false)
	tx = tx.Select("profiles.id as profile_id", "profiles.friendly_name as profile_friendly_name", "profiles.datasheet_id as id", "files.type", "files.created_at", "files.created_by", "files.updated_at", "files.updated_by", "files.deleted_at", "files.deleted_by")

	tx.Joins(fmt.Sprintf("left join %s on %s.profile_id = %s.id", profileTableName, tableName, profileTableName))
	tx.Joins(fmt.Sprintf("left join %s on %s.datasheet_id = %s.id", fileTableName, profileTableName, fileTableName))
	tx.Scopes(
		apertureDao.ByQuoteId(daoOperator.Equals, *id),
		apertureDao.ByProfileId(daoOperator.IsNot, nil),
		apertureDao.ByDatasheetId(daoOperator.IsNot, nil),
	)
	tx.Group("profiles.id, files.type, files.created_at, files.created_by, files.updated_at, files.updated_by, files.deleted_at, files.deleted_by")

	// The datasheets of the aperture types (profiles) go first
	itemListMap := []map[string]any{}
	if err := tx.Find(&itemListMap).Error; err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch items: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Followed by the datasheets of the glass types, fetched the same way
	glassTableName := glassDao.GetInstance().TableName()

	glassTx := apertureDao.GetBasicTransaction(ctx, apertureDao.NewEntity, false)
	glassTx = glassTx.Select("glasses.id as glass_id", "glasses.friendly_name as glass_friendly_name", "glasses.datasheet_id as id", "files.type", "files.created_at", "files.created_by", "files.updated_at", "files.updated_by", "files.deleted_at", "files.deleted_by")

	glassTx.Joins(fmt.Sprintf("left join %s on %s.glass_id = %s.id", glassTableName, tableName, glassTableName))
	glassTx.Joins(fmt.Sprintf("left join %s on %s.datasheet_id = %s.id", fileTableName, glassTableName, fileTableName))
	glassTx.Scopes(apertureDao.ByQuoteId(daoOperator.Equals, *id))
	glassTx.Where(fmt.Sprintf("%s.glass_id IS NOT NULL", tableName))
	glassTx.Where(fmt.Sprintf("%s.datasheet_id IS NOT NULL", glassTableName))
	glassTx.Group("glasses.id, files.type, files.created_at, files.created_by, files.updated_at, files.updated_by, files.deleted_at, files.deleted_by")
	glassTx.Order("glasses.id")

	glassItemListMap := []map[string]any{}
	if err := glassTx.Find(&glassItemListMap).Error; err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch items: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	itemListMap = append(itemListMap, glassItemListMap...)

	pagination := daoQuery.NewPagination()
	pagination.Data = itemListMap
	pagination.Meta.TotalItems = len(itemListMap)
	pagination.Meta.Limit = int(pagination.Meta.TotalItems)
	pagination.Meta.TotalPages = 1

	c.JSON(http.StatusOK, pagination)
}

func GetCalculationsById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	id, err := apiUtils.HttpPathUInt64(c, "quote_id")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	result, err := quoteDao.GetInstance().ReadCalculationsByID(ctx, *id)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch item by id: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, result)
}

func GenerateClientInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	quoteId, err := apiUtils.HttpPathUInt64(c, "quote_id")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	includeBom, err := apiUtils.HttpQueryBool(c, "include-bom", true)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	selector := daoQuery.Selector{
		SubItems: []string{"clients"},
	}
	quoteI, _, err := quoteDao.GetInstance().ReadByID(ctx, *quoteId, &selector)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}
	quote := quoteI.(*quoteEntity.QuoteEntity)

	calcMap, err := quoteDao.GetInstance().ReadCalculationsByID(ctx, *quoteId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}

	quote.TotalApertures = calcMap["total_apertures"].(int)
	quote.NetsQuantity = calcMap["nets_quantity"].(int)
	quote.Discount = calcMap["discount"].(float64)
	quote.SubTotal = calcMap["sub_total"].(float64)
	quote.Vat = calcMap["vat"].(float64)
	quote.GrossProfit = calcMap["gross_profit"].(float64)
	quote.Total = calcMap["total"].(float64)
	quote.TotalArea = calcMap["total_area"].(float64)
	quote.TotalCost = calcMap["total_cost"].(float64)

	var mainApertureList []*apertureEntity.ApertureEntity
	// Get the list of apertures
	{
		noImageBytes, err := imageUtils.FromFilePath(pathUtils.ToImagePath("no-image.jpg"))
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
		err = apertureDao.GetInstance().ReadMainByQuoteId(ctx, *quoteId, p, &selector)
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
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
						apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
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
			apertureDao.ByQuoteId(daoOperator.Equals, *quoteId),
			apertureDao.ByProfileId(daoOperator.IsNot, nil),
			apertureDao.ByDatasheetId(daoOperator.IsNot, nil),
		)
		tx.Group("profiles.id, files.type")

		itemListMap := []map[string]any{}
		if err := tx.Find(&itemListMap).Error; err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
			return
		}

		for _, item := range itemListMap {
			id := item["id"].(uint64)
			fileI, found, err := fileDao.ReadByID(ctx, id, nil)
			if err != nil {
				apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
				return
			}

			if found {
				file := fileI.(*fileEntity.FileEntity)
				datasheetList = append(datasheetList, file.Data)
			}
		}
	}

	pdfBytes, err := clientPdf.Generate(quote, mainApertureList, datasheetList, *includeBom)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=Client Invoice")
	c.Header("Content-Type", "application/pdf")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Writer.Write(pdfBytes)
}

func GenerateSupplierInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	quoteId, err := apiUtils.HttpPathUInt64(c, "quote_id")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	selector := daoQuery.Selector{
		SubItems: []string{"clients"},
	}
	quoteI, _, err := quoteDao.GetInstance().ReadByID(ctx, *quoteId, &selector)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}
	quote := quoteI.(*quoteEntity.QuoteEntity)

	var mainApertureList []*apertureEntity.ApertureEntity
	// Get the list of apertures
	{
		noImageBytes, err := imageUtils.FromFilePath(pathUtils.ToImagePath("no-image.jpg"))
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
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

		allApertures := []*apertureEntity.ApertureEntity{}
		apertureMap := map[uint64]*apertureEntity.ApertureEntity{}
		orderedKeys := []uint64{}

		// The main apertures
		{
			err = apertureDao.GetInstance().ReadMainByQuoteId(ctx, *quoteId, p, &selector)
			if err != nil {
				apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
				return
			}

			for _, apertureI := range p.Data.([]any) {
				aperture := apertureI.(*apertureEntity.ApertureEntity)

				allApertures = append(allApertures, aperture)
				apertureMap[aperture.ID] = aperture
				orderedKeys = append(orderedKeys, aperture.ID)
			}
		}

		// The child apertures
		{
			p.Data = []any{}
			err = apertureDao.GetInstance().ReadChildrenByQuoteId(ctx, *quoteId, p, &selector)
			if err != nil {
				apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
				return
			}

			for _, apertureI := range p.Data.([]any) {
				aperture := apertureI.(*apertureEntity.ApertureEntity)

				allApertures = append(allApertures, aperture)
				parentAperture := apertureMap[*aperture.ParentApertureID]
				parentAperture.ChildApertures = append(parentAperture.ChildApertures, aperture)
			}
		}

		for _, key := range orderedKeys {
			mainApertureList = append(mainApertureList, apertureMap[key])
		}

		// Fill the apertures' image
		{
			for _, aperture := range allApertures {
				aperture.Image = noImageEntity
				if aperture.ImageId != nil {
					imgEntityI, found, err := fileDao.GetInstance().ReadByID(ctx, *aperture.ImageId, nil)
					if err != nil {
						apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
						return
					}

					if found {
						imgEntity := imgEntityI.(*fileEntity.FileEntity)
						aperture.Image = imgEntity
					}
				}
			}
		}
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
			apertureDao.ByQuoteId(daoOperator.Equals, *quoteId),
			apertureDao.ByProfileId(daoOperator.IsNot, nil),
			apertureDao.ByDatasheetId(daoOperator.IsNot, nil),
		)
		tx.Group("profiles.id, files.type")

		itemListMap := []map[string]any{}
		if err := tx.Find(&itemListMap).Error; err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
			return
		}

		for _, item := range itemListMap {
			id := item["id"].(uint64)
			fileI, found, err := fileDao.ReadByID(ctx, id, nil)
			if err != nil {
				apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
				return
			}

			if found {
				file := fileI.(*fileEntity.FileEntity)
				datasheetList = append(datasheetList, file.Data)
			}
		}
	}

	// The terms and conditions PDF, if there is one, goes at the end of the report
	termsAndConditions, err := termsAndConditionsService.GetPdf(ctx)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch the terms and conditions: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	pdfBytes, err := supplierPdf.Generate(quote, mainApertureList, datasheetList, termsAndConditions)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=Client Invoice")
	c.Header("Content-Type", "application/pdf")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Writer.Write(pdfBytes)
}
