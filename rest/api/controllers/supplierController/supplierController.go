package supplierController

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daos/supplierDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/supplierEntity"
	"Quote/logger"
	"Quote/model/baseModel"
	"Quote/model/models/supplierModel"
	"Quote/rest/api/apiUtils"
	"Quote/rest/api/baseController"
	"net/http"

	"github.com/gin-gonic/gin"
)

const withAuthorisation = true

func GetNewEntity() baseEntity.BaseEntityI {
	return supplierEntity.New()
}
func GetNewModel() baseModel.BaseModelI {
	return supplierModel.New()
}
func GetDao() baseDao.BaseDaoI {
	return supplierDao.GetInstance()
}
func GetSystemFields() []string {
	return baseController.GetSystemFields()
}

func Create(c *gin.Context) {
	baseController.Create(c, GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func Update(c *gin.Context) {
	baseController.Update(c, "supplier_id", GetNewModel(), GetDao(), GetSystemFields(), withAuthorisation)
}

func SoftDelete(c *gin.Context) {
	baseController.SoftDelete(c, "supplier_id", GetDao(), withAuthorisation)
}

func GetById(c *gin.Context) {
	baseController.GetById(c, "supplier_id", GetDao())
}

func Get(c *gin.Context) {
	pagination := baseController.Get(c, GetDao())
	if pagination == nil {
		apiUtils.HttpError(c.Writer, logger.GetInstance(), "Failed to fetch suppliers", http.StatusInternalServerError)
		return
	}

	itemsMap := []map[string]any{}
	for _, item := range pagination.Data.([]any) {
		currentItem := item.(*supplierEntity.SupplierEntity)
		itemsMap = append(itemsMap, currentItem.ToJson())
	}
	pagination.Data = itemsMap

	c.JSON(http.StatusOK, pagination)
}
