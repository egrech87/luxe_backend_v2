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
	"Quote/service/quoteService"
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

// Duplicate creates a copy of a quote together with all its apertures and their child apertures.
// The copy gets the next quote number. It responds with the new quote.
func Duplicate(c *gin.Context) {
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

	newId, found, err := quoteService.Duplicate(ctx, *id, authUserId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to duplicate item: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "item not available", http.StatusNotFound)
		return
	}

	// Return the new quote, with its calculated fields, in the same shape as GET /quotes/by/id/:quote_id
	selector, err := apiUtils.HttpHeaderSelector(c)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}
	copiedQuote, found, err := GetDao().ReadByID(ctx, newId, selector)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch duplicated item: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "duplicated item not available", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, copiedQuote.ToJson())
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

	// The datasheets of the aperture types (profiles), followed by those of the glass types
	datasheetList, err := getDatasheetPdfs(ctx, *quoteId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}

	// The terms and conditions PDF, if there is one, goes at the very end of the report
	termsAndConditions, err := termsAndConditionsService.GetPdf(ctx)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch the terms and conditions: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	pdfBytes, err := clientPdf.Generate(quote, mainApertureList, datasheetList, termsAndConditions, *includeBom)
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

	// The datasheets of the aperture types (profiles), followed by those of the glass types
	datasheetList, err := getDatasheetPdfs(ctx, *quoteId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
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

// getDatasheetPdfs returns the datasheet PDFs for a quote: first those of the aperture types (profiles) used by its apertures,
// then those of the glass types, each fetched the same way and ordered by type id.
func getDatasheetPdfs(ctx context.Context, quoteId uint64) ([][]byte, error) {
	apertureDao := apertureDao.GetInstance()
	fileDao := fileDao.GetInstance()
	tableName := apertureDao.TableName()
	fileTableName := fileDao.TableName()

	sources := []struct {
		typeTableName string
		typeIdColumn  string
	}{
		{profileDao.GetInstance().TableName(), "profile_id"},
		{glassDao.GetInstance().TableName(), "glass_id"},
	}

	datasheetList := [][]byte{}
	for _, source := range sources {
		tx := apertureDao.GetBasicTransaction(ctx, apertureDao.NewEntity, false)
		tx = tx.Select(fmt.Sprintf("%s.datasheet_id as id", source.typeTableName))

		tx.Joins(fmt.Sprintf("left join %s on %s.%s = %s.id", source.typeTableName, tableName, source.typeIdColumn, source.typeTableName))
		tx.Joins(fmt.Sprintf("left join %s on %s.datasheet_id = %s.id", fileTableName, source.typeTableName, fileTableName))
		tx.Where(fmt.Sprintf("%s.type = ?", fileTableName), fileTypes.Pdf)
		tx.Where(fmt.Sprintf("%s.%s IS NOT NULL", tableName, source.typeIdColumn))
		tx.Where(fmt.Sprintf("%s.datasheet_id IS NOT NULL", source.typeTableName))
		tx.Scopes(apertureDao.ByQuoteId(daoOperator.Equals, quoteId))
		tx.Group(fmt.Sprintf("%s.id, %s.type", source.typeTableName, fileTableName))
		tx.Order(fmt.Sprintf("%s.id", source.typeTableName))

		itemListMap := []map[string]any{}
		if err := tx.Find(&itemListMap).Error; err != nil {
			return nil, err
		}

		for _, item := range itemListMap {
			id := item["id"].(uint64)
			fileI, found, err := fileDao.ReadByID(ctx, id, nil)
			if err != nil {
				return nil, err
			}

			if found {
				datasheetList = append(datasheetList, fileI.(*fileEntity.FileEntity).Data)
			}
		}
	}

	return datasheetList, nil
}
