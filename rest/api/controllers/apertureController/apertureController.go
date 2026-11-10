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

// maxDuplicateAmount is the largest number of copies a single Duplicate request may create.
const maxDuplicateAmount = 100

// Duplicate creates `amount` persisted copies (default 1, from the "amount" query parameter) of an aperture and of all its child apertures.
// It responds with only the first copy that was created.
func Duplicate(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	authUserId, err := baseController.GetAuthenticatedId(c)
	if withAuthorisation && err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnauthorized)
		return
	}

	id, err := apiUtils.HttpPathUInt64(c, "aperture_id")
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	amount := uint64(1)
	{
		amountParam, err := apiUtils.HttpQueryUInt64(c, "amount", false)
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
			return
		}
		if amountParam != nil {
			amount = *amountParam
		}
		if amount < 1 || amount > maxDuplicateAmount {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("'amount' query parameter must be between 1 and %d", maxDuplicateAmount), http.StatusBadRequest)
			return
		}
	}

	dao := apertureDao.GetInstance()

	newId, found, err := dao.Duplicate(ctx, *id, int(amount), authUserId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to duplicate item: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "item not available", http.StatusNotFound)
		return
	}

	// Return the first copy, with its calculated fields, in the same shape as GET /apertures/by/id/:aperture_id
	selector, err := apiUtils.HttpHeaderSelector(c)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}
	copiedAperture, found, err := dao.ReadByID(ctx, newId, selector)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch duplicated item: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "duplicated item not available", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, copiedAperture.ToJson())
}
