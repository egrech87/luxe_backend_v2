package baseController

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/logger"
	"Quote/model/baseModel"
	"Quote/rest/api/apiUtils"
	"Quote/utils/mapUtils"
	"context"
	"fmt"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	Timeout time.Duration = 30 * time.Second
)

func CleanInput(json map[string]any, systemFields []string) {
	for _, sysField := range systemFields {
		if mapUtils.HasKey(json, sysField) {
			delete(json, sysField)
		}
	}
}

func GetSystemFields() []string {
	return []string{"id", "created_at", "created_by", "updated_at", "updated_by", "deleted", "deleted_at", "deleted_by"}
}

func GetAuthenticatedId(c *gin.Context) (uint64, error) {
	username, _, ok := c.Request.BasicAuth()
	if !ok {
		return 0, errors.Errorf("Failed to get Authentication details")
	}

	return strconv.ParseUint(username, 10, 64)
}

func Create(c *gin.Context, model baseModel.BaseModelI, dao baseDao.BaseDaoI, systemFields []string, withAuthorisation bool) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	authUserId, err := GetAuthenticatedId(c)
	if withAuthorisation && err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnauthorized)
		return
	}

	// Get the model from the request
	{
		inputJson, err := apiUtils.HttpBodyJson(c, true)
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
			return
		}
		CleanInput(inputJson, systemFields)
		model.FillFromJson(inputJson)
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

func Update(c *gin.Context, idQueryParameter string, originalModel baseModel.BaseModelI, dao baseDao.BaseDaoI, systemFields []string, withAuthorisation bool) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	authUserId, err := GetAuthenticatedId(c)
	if withAuthorisation && err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnauthorized)
		return
	}

	id, err := apiUtils.HttpPathUInt64(c, idQueryParameter)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	inputJson, err := apiUtils.HttpBodyJson(c, true)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}
	CleanInput(inputJson, GetSystemFields())

	// Get the model from the database
	{
		originalEntity, found, err := dao.ReadByID(ctx, *id, nil)
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
			return
		}
		if !found {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("item %d not found", *id), http.StatusNotFound)
			return
		}
		originalModel.SetEntity(originalEntity)
	}

	targetModel, err := originalModel.Clone()
	// Fill in the target model with the request
	{
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to clone item: %s", err.Error()), http.StatusUnprocessableEntity)
			return
		}
		if err = targetModel.FillFromJson(inputJson); err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fill item from input: %s", err.Error()), http.StatusUnprocessableEntity)
			return
		}
	}

	updates := targetModel.GetUpdates(originalModel)

	if len(updates) > 0 {
		err = targetModel.Validate(ctx)
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnprocessableEntity)
			return
		}

		err = dao.Update(ctx, *id, targetModel.GetEntity(), mapUtils.GetKeys(updates), authUserId)
		if err != nil {
			apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to update item: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

func SoftDelete(c *gin.Context, idQueryParameter string, dao baseDao.BaseDaoI, withAuthorisation bool) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	authUserId, err := GetAuthenticatedId(c)
	if withAuthorisation && err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnauthorized)
		return
	}

	id, err := apiUtils.HttpPathUInt64(c, idQueryParameter)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	_, err = dao.SoftDelete(ctx, *id, authUserId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to delete the item: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func GetById(c *gin.Context, idParamName string, dao baseDao.BaseDaoI) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	selector, err := apiUtils.HttpHeaderSelector(c)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	id, err := apiUtils.HttpPathUInt64(c, idParamName)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	item, found, err := dao.ReadByID(ctx, *id, selector)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch item by id: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "item not available", http.StatusNotFound)
		return
	}

	result := item.ToJson()

	c.JSON(http.StatusOK, result)
}

func Get(c *gin.Context, dao baseDao.BaseDaoI) *daoQuery.Pagination {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	selector, err := apiUtils.HttpHeaderSelector(c)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return nil
	}

	pagination, err := apiUtils.HttpQueryPagination(c)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return nil
	}
	if pagination.Meta.Sort == "" {
		pagination.Meta.Sort = fmt.Sprintf("%s.created_at DESC", dao.TableName())
	}

	err = dao.Read(ctx, pagination, selector)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to fetch items: %s", err.Error()), http.StatusInternalServerError)
		return nil
	}

	return pagination
}

// Duplicate creates a persisted copy of the item identified by the path parameter and responds with the new item's id.
func Duplicate(c *gin.Context, idQueryParameter string, dao baseDao.DuplicableI, withAuthorisation bool) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	authUserId, err := GetAuthenticatedId(c)
	if withAuthorisation && err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnauthorized)
		return
	}

	id, err := apiUtils.HttpPathUInt64(c, idQueryParameter)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusBadRequest)
		return
	}

	newId, found, err := dao.Duplicate(ctx, *id, authUserId)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("failed to duplicate item: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), fmt.Sprintf("item %d not found", *id), http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, newId)
}
