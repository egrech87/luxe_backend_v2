package userController

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daos/userDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/userEntity"
	"Quote/logger"
	"Quote/model/baseModel"
	"Quote/model/models/userModel"
	"Quote/rest/api/apiUtils"
	"Quote/rest/api/baseController"
	"Quote/utils/timeUtils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

const withAuthorisation = false

func GetNewEntity() baseEntity.BaseEntityI {
	return userEntity.New()
}
func GetNewModel() baseModel.BaseModelI {
	return userModel.New()
}
func GetDao() baseDao.BaseDaoI {
	return userDao.GetInstance()
}
func GetSystemFields() []string {
	return append(baseController.GetSystemFields(), "last_login")
}

func Create(c *gin.Context) {
	baseController.Create(c, GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func Update(c *gin.Context) {
	baseController.Update(c, "user_id", GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func SoftDelete(c *gin.Context) {
	baseController.SoftDelete(c, "user_id", GetDao(), withAuthorisation)
}

func GetById(c *gin.Context) {
	baseController.GetById(c, "user_id", GetDao())
}

func Get(c *gin.Context) {
	pagination := baseController.Get(c, GetDao())
	if pagination == nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	itemsMap := []map[string]any{}
	for _, item := range pagination.Data.([]any) {
		currentItem := item.(*userEntity.UserEntity)
		itemsMap = append(itemsMap, currentItem.ToJson())
	}
	pagination.Data = itemsMap

	c.JSON(http.StatusOK, pagination)
}

func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), baseController.Timeout)
	defer cancel()

	authUserId, err := baseController.GetAuthenticatedId(c)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusUnauthorized)
		return
	}

	userEntityI, found, err := userDao.GetInstance().ReadByID(ctx, authUserId, nil)
	if err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "User not found", http.StatusNotFound)
		return
	}

	userEntity := userEntityI.(*userEntity.UserEntity)
	userEntity.SetLastLogin(timeUtils.CurrentLocalisedDateTime())
	if err = userDao.GetInstance().Update(ctx, authUserId, userEntity, []string{"last_login"}, 0); err != nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), err.Error(), http.StatusInternalServerError)
		return
	}

	result := map[string]any{
		"item": userEntity.ToJson(),
	}

	c.JSON(http.StatusOK, result)
}
