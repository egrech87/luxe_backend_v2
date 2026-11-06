package clientController

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daos/clientDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/clientEntity"
	"Quote/logger"
	"Quote/model/baseModel"
	"Quote/model/models/clientModel"
	"Quote/rest/api/apiUtils"
	"Quote/rest/api/baseController"
	"net/http"

	"github.com/gin-gonic/gin"
)

const withAuthorisation = true

func GetNewEntity() baseEntity.BaseEntityI {
	return clientEntity.New()
}
func GetNewModel() baseModel.BaseModelI {
	return clientModel.New()
}
func GetDao() baseDao.BaseDaoI {
	return clientDao.GetInstance()
}
func GetSystemFields() []string {
	return baseController.GetSystemFields()
}

func Create(c *gin.Context) {
	baseController.Create(c, GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func Update(c *gin.Context) {
	baseController.Update(c, "client_id", GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func SoftDelete(c *gin.Context) {
	baseController.SoftDelete(c, "client_id", GetDao(), withAuthorisation)
}

func GetById(c *gin.Context) {
	baseController.GetById(c, "client_id", GetDao())
}

func Get(c *gin.Context) {
	pagination := baseController.Get(c, GetDao())
	if pagination == nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "Failed to fetch clients", http.StatusInternalServerError)
		return
	}

	itemsMap := []map[string]any{}
	for _, item := range pagination.Data.([]any) {
		currentItem := item.(*clientEntity.ClientEntity)
		itemsMap = append(itemsMap, currentItem.ToJson())
	}
	pagination.Data = itemsMap

	c.JSON(http.StatusOK, pagination)
}
