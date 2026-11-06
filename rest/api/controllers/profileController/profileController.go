package profileController

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daos/profileDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/profileEntity"
	"Quote/logger"
	"Quote/model/baseModel"
	"Quote/model/models/profileModel"
	"Quote/rest/api/apiUtils"
	"Quote/rest/api/baseController"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const withAuthorisation = true

func GetNewEntity() baseEntity.BaseEntityI {
	return profileEntity.New()
}
func GetNewModel() baseModel.BaseModelI {
	return profileModel.New()
}
func GetDao() baseDao.BaseDaoI {
	return profileDao.GetInstance()
}
func GetSystemFields() []string {
	return baseController.GetSystemFields()
}

func Create(c *gin.Context) {
	baseController.Create(c, GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func Update(c *gin.Context) {
	baseController.Update(c, "profile_id", GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func SoftDelete(c *gin.Context) {
	baseController.SoftDelete(c, "profile_id", GetDao(), withAuthorisation)
}

func GetById(c *gin.Context) {
	baseController.GetById(c, "profile_id", GetDao())
}

func Get(c *gin.Context) {
	pagination := baseController.Get(c, GetDao())
	if pagination == nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "Failed to fetch profiles", http.StatusInternalServerError)
		return
	}

	itemsMap := []map[string]any{}
	for _, item := range pagination.Data.([]any) {
		currentItem := item.(*profileEntity.ProfileEntity)
		itemsMap = append(itemsMap, currentItem.ToJson())
	}
	pagination.Data = itemsMap

	c.JSON(http.StatusOK, pagination)
}

func GetMaxDimensionsInPriceMatrix(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	id, err := apiUtils.HttpPathUInt64(c, "profile_id")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	width, err := apiUtils.HttpPathInt(c, "width")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	item, found, err := profileDao.GetInstance().ReadByID(ctx, *id, nil)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch item by id: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "item not available", http.StatusNotFound)
		return
	}

	profile := item.(*profileEntity.ProfileEntity)

	var priceMatrixMap map[string]any
	err = json.Unmarshal([]byte(profile.PriceMatrix), &priceMatrixMap)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "failed to convert the price matrix", http.StatusInternalServerError)
		return
	}

	chosenWidth := math.MaxInt
	maxWidth := 0
	for key := range priceMatrixMap {
		currentWidth, err := strconv.Atoi(key)
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), "failed to convert the price matrix's width", http.StatusInternalServerError)
			return
		}

		if currentWidth >= *width && currentWidth < chosenWidth {
			chosenWidth = currentWidth
		}

		if currentWidth > maxWidth {
			maxWidth = currentWidth
		}
	}

	if chosenWidth == math.MaxInt {
		chosenWidth = maxWidth
	}

	widthStr := fmt.Sprintf("%d", chosenWidth)
	widthMap := priceMatrixMap[widthStr].(map[string]any)

	maxHeight := 0
	for key := range widthMap {
		currentHeight, err := strconv.Atoi(key)
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), "failed to convert the price matrix's height", http.StatusInternalServerError)
			return
		}

		if currentHeight > maxHeight {
			maxHeight = currentHeight
		}
	}

	result := map[string]int{
		"max_width":  maxWidth,
		"max_height": maxHeight,
	}

	c.JSON(http.StatusOK, result)
}
