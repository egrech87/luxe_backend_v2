package glassController

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daos/glassDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/glassEntity"
	"Quote/logger"
	"Quote/model/baseModel"
	"Quote/model/models/glassModel"
	"Quote/rest/api/apiUtils"
	"Quote/rest/api/baseController"
	"net/http"

	"github.com/gin-gonic/gin"
)

const withAuthorisation = true

func GetNewEntity() baseEntity.BaseEntityI {
	return glassEntity.New()
}
func GetNewModel() baseModel.BaseModelI {
	return glassModel.New()
}
func GetDao() baseDao.BaseDaoI {
	return glassDao.GetInstance()
}
func GetSystemFields() []string {
	return baseController.GetSystemFields()
}

func Create(c *gin.Context) {
	baseController.Create(c, GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func Update(c *gin.Context) {
	baseController.Update(c, "glass_id", GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func Duplicate(c *gin.Context) {
	baseController.Duplicate(c, "glass_id", glassDao.GetInstance(), withAuthorisation)
}

func SoftDelete(c *gin.Context) {
	baseController.SoftDelete(c, "glass_id", GetDao(), withAuthorisation)
}

func GetById(c *gin.Context) {
	baseController.GetById(c, "glass_id", GetDao())
}

func Get(c *gin.Context) {
	pagination := baseController.Get(c, GetDao())
	if pagination == nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "Failed to fetch glasses", http.StatusInternalServerError)
		return
	}

	itemsMap := []map[string]any{}
	for _, item := range pagination.Data.([]any) {
		currentItem := item.(*glassEntity.GlassEntity)
		itemsMap = append(itemsMap, currentItem.ToJson())
	}
	pagination.Data = itemsMap

	c.JSON(http.StatusOK, pagination)
}
