package apertureController

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daos/apertureDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/apertureEntity"
	"Quote/logger"
	"Quote/model/baseModel"
	"Quote/model/models/apertureModel"
	"Quote/rest/api/apiUtils"
	"Quote/rest/api/baseController"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const withAuthorisation = true

func GetNewEntity() baseEntity.BaseEntityI {
	return apertureEntity.New()
}
func GetNewModel() baseModel.BaseModelI {
	return apertureModel.New()
}
func GetDao() baseDao.BaseDaoI {
	return apertureDao.GetInstance()
}
func GetSystemFields() []string {
	return append(baseController.GetSystemFields(), "is_combined")
}

func Create(c *gin.Context) {
	baseController.Create(c, GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func Update(c *gin.Context) {
	baseController.Update(c, "aperture_id", GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func SoftDelete(c *gin.Context) {
	baseController.SoftDelete(c, "aperture_id", GetDao(), withAuthorisation)
}

func GetById(c *gin.Context) {
	baseController.GetById(c, "aperture_id", GetDao())
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
	selector.ExcludeFields = append(selector.ExcludeFields, "type")
	selector.ExcludeFields = append(selector.ExcludeFields, "handle")
	selector.ExcludeFields = append(selector.ExcludeFields, "area")
	selector.ExcludeFields = append(selector.ExcludeFields, "overall_u_value")
	selector.ExcludeFields = append(selector.ExcludeFields, "color_price_factor")
	selector.ExcludeFields = append(selector.ExcludeFields, "glass_price_factor")
	selector.ExcludeFields = append(selector.ExcludeFields, "color_id")
	selector.ExcludeFields = append(selector.ExcludeFields, "glass_id")
	selector.ExcludeFields = append(selector.ExcludeFields, "matrix_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "nets_quantity")
	selector.ExcludeFields = append(selector.ExcludeFields, "nets_price")
	selector.ExcludeFields = append(selector.ExcludeFields, "nets_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "color_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "glass_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "other_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "other_cost_comment")
	selector.ExcludeFields = append(selector.ExcludeFields, "supplier_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "shipping_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "installation_per_unit")
	selector.ExcludeFields = append(selector.ExcludeFields, "installation_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "contingency_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "item_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "profit_margin")
	selector.ExcludeFields = append(selector.ExcludeFields, "total_area")
	selector.ExcludeFields = append(selector.ExcludeFields, "total_cost")
	selector.ExcludeFields = append(selector.ExcludeFields, "total_profit")
	selector.ExcludeFields = append(selector.ExcludeFields, "updated_at")

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
		currentItem := item.(*apertureEntity.ApertureEntity)
		itemsMap = append(itemsMap, currentItem.ToJson())
	}
	pagination.Data = itemsMap

	c.JSON(http.StatusOK, pagination)
}

func Copy(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	dao := GetDao()

	selector, err := apiUtils.HttpHeaderSelector(c)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	id, err := apiUtils.HttpPathUInt64(c, "aperture_id")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	apertureI, found, err := dao.ReadByID(ctx, *id, selector)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch items: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch item by id: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "item not available", http.StatusNotFound)
		return
	}

	copiedAperture, err := apertureI.Clone()
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}

	copiedAperture.SetID(0)

	c.JSON(http.StatusOK, copiedAperture.ToJson())
}
