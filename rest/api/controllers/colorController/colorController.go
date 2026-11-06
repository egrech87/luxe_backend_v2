package colorController

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daos/colorDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/colorEntity"
	"Quote/logger"
	"Quote/model/baseModel"
	"Quote/model/models/colorModel"
	"Quote/rest/api/apiUtils"
	"Quote/rest/api/baseController"
	"net/http"

	"github.com/gin-gonic/gin"
)

const withAuthorisation = true

func GetNewEntity() baseEntity.BaseEntityI {
	return colorEntity.New()
}
func GetNewModel() baseModel.BaseModelI {
	return colorModel.New()
}
func GetDao() baseDao.BaseDaoI {
	return colorDao.GetInstance()
}
func GetSystemFields() []string {
	return baseController.GetSystemFields()
}

func Create(c *gin.Context) {
	baseController.Create(c, GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func Update(c *gin.Context) {
	baseController.Update(c, "color_id", GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func SoftDelete(c *gin.Context) {
	baseController.SoftDelete(c, "color_id", GetDao(), withAuthorisation)
}

func GetById(c *gin.Context) {
	baseController.GetById(c, "color_id", GetDao())
}

func Get(c *gin.Context) {
	pagination := baseController.Get(c, GetDao())
	if pagination == nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "Failed to fetch colors", http.StatusInternalServerError)
		return
	}

	itemsMap := []map[string]any{}
	for _, item := range pagination.Data.([]any) {
		currentItem := item.(*colorEntity.ColorEntity)
		itemsMap = append(itemsMap, currentItem.ToJson())
	}
	pagination.Data = itemsMap

	c.JSON(http.StatusOK, pagination)
}
