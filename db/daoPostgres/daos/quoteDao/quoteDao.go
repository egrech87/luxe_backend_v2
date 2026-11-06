package quoteDao

import (
	"Quote/db"
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/daoPostgres/daoQuery/enums/daoOperator"
	"Quote/db/daoPostgres/daos/clientDao"
	"Quote/db/daoPostgres/daos/userDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/quoteEntity"
	"Quote/utils/anyUtils"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type QuoteDao struct {
	baseDao.BaseDao
}

var instance *QuoteDao

func GetInstance() *QuoteDao {
	if instance == nil {
		instance = new(QuoteDao)
	}

	return instance
}

func (d *QuoteDao) TableName() string {
	return d.NewEntity().TableName()
}
func (d *QuoteDao) NewEntity() baseEntity.BaseEntityI {
	return quoteEntity.New()
}

func (d *QuoteDao) ApplySelect(tx *gorm.DB, subItemList []string, excludeFields []string) (*gorm.DB, error) {
	fieldNames := d.FieldNamesForTransaction(d, nil)
	tableName := d.TableName()

	for _, subItem := range subItemList {
		switch subItem {
		case "clients":
			clientDao := clientDao.GetInstance()
			clientTableName := clientDao.TableName()
			clientFieldNames := clientDao.FieldNamesForTransaction(clientDao, nil)
			fieldNames = append(fieldNames, clientFieldNames...)
			tx.Joins(fmt.Sprintf("left join %s on %s.client_id = %s.id", clientTableName, tableName, clientTableName))
		case "createuser":
			userDao := userDao.GetInstance()
			userTableName := userDao.TableName()
			userTableAlias := subItem
			userFieldNames := userDao.FieldNamesForTransaction(userDao, &userTableAlias)
			fieldNames = append(fieldNames, userFieldNames...)
			tx.Joins(fmt.Sprintf("left join %s as %s on %s.created_by = %s.id", userTableName, userTableAlias, tableName, userTableAlias))
		case "updateuser":
			userDao := userDao.GetInstance()
			userTableName := userDao.TableName()
			userTableAlias := subItem
			userFieldNames := userDao.FieldNamesForTransaction(userDao, &userTableAlias)
			fieldNames = append(fieldNames, userFieldNames...)
			tx.Joins(fmt.Sprintf("left join %s as %s on quote_updated_by(%s.id) = %s.id", userTableName, userTableAlias, tableName, userTableAlias))
		}
	}

	fieldNames = append(fieldNames, fmt.Sprintf("quote_updated_at(quotes.id) as %s_updated_at", tableName))

	tx.Select(fieldNames)

	return tx, nil
}
func (d *QuoteDao) ApplyFilters(tx *gorm.DB, filters []daoQuery.Filter) (*gorm.DB, error) {
	for _, filter := range filters {
		switch filter.Field {
		case "friendly_name":
			tx.Scopes(d.ByFriendlyName(filter.Operator, filter.Value.(string)))
		case "client.name":
			tx.Scopes(d.ByClientName(filter.Operator, filter.Value.(string)))
		case "client.surname":
			tx.Scopes(d.ByClientSurname(filter.Operator, filter.Value.(string)))
		case "locality":
			tx.Scopes(d.ByLocality(filter.Operator, filter.Value.(string)))
		case "country":
			tx.Scopes(d.ByCountry(filter.Operator, filter.Value.(string)))
		}
	}
	return tx, nil
}

// #region Basic CRUD operations

func (d *QuoteDao) Create(ctx context.Context, quote baseEntity.BaseEntityI, createdBy uint64) error {
	return d.BaseDao.Create(ctx, quote, createdBy)
}

func (d *QuoteDao) Update(ctx context.Context, id uint64, quote baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error {
	return d.BaseDao.Update(ctx, d.TableName(), id, quote, fieldsToUpdate, updatedBy)
}

func (d *QuoteDao) SoftDelete(ctx context.Context, id uint64, deletedBy uint64) (deletedAt time.Time, err error) {
	return d.BaseDao.SoftDelete(ctx, id, quoteEntity.TableName(), deletedBy)
}

func (d *QuoteDao) ReadByID(ctx context.Context, id uint64, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	return d.BaseDao.ReadByID(ctx, d, id, selector)
}

func (d *QuoteDao) Read(ctx context.Context, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false), d, pagination, selector)
}

// #endregion

func (d *QuoteDao) ReadNextNumber(ctx context.Context) (number int, err error) {
	resultMap := map[string]any{}
	result := db.GetDB().WithContext(ctx).Raw("Select quote_next_number() as next_number").Find(&resultMap)
	if result.Error != nil {
		return 0, result.Error
	}

	nextNumber, err := anyUtils.ToInt(resultMap["next_number"])
	if err != nil {
		return 0, err
	}

	return nextNumber, nil
}

func (d *QuoteDao) ReadCalculationsByID(ctx context.Context, quoteId uint64) (resultMap map[string]any, err error) {
	fieldNames := []string{}

	fieldNames = append(fieldNames, "total_apertures_in_quote(quotes.id) as total_apertures")
	fieldNames = append(fieldNames, "quote_nets_quantity(quotes.id) as nets_quantity")
	fieldNames = append(fieldNames, "quote_discount(quotes.id) as discount")
	fieldNames = append(fieldNames, "quote_gross_profit(quotes.id) as gross_profit")
	fieldNames = append(fieldNames, "quote_sub_total(quotes.id) as sub_total")
	fieldNames = append(fieldNames, "quote_total(quotes.id) as total")
	fieldNames = append(fieldNames, "quote_vat(quotes.id) as vat")
	fieldNames = append(fieldNames, "quote_total_area(quotes.id) as total_area")
	fieldNames = append(fieldNames, "quote_total_cost(quotes.id) as total_cost")

	resultMap = map[string]any{}
	result := db.GetDB().WithContext(ctx).Table(d.TableName()).Select(fieldNames).Scopes(d.ById(d.TableName(), quoteId)).Find(&resultMap)
	if result.Error != nil {
		return resultMap, result.Error
	}

	entity := quoteEntity.New()
	entity.FillFromJson(resultMap)

	resultMap["total_apertures"] = entity.TotalApertures
	resultMap["nets_quantity"] = entity.NetsQuantity
	resultMap["discount"] = entity.Discount
	resultMap["gross_profit"] = entity.GrossProfit
	resultMap["sub_total"] = entity.SubTotal
	resultMap["total"] = entity.Total
	resultMap["vat"] = entity.Vat
	resultMap["total_area"] = entity.TotalArea
	resultMap["total_cost"] = entity.TotalCost
	resultMap["updated_at"] = entity.UpdatedAt

	return resultMap, nil
}

// #region Scopes

// #region WHERE clause

func (d *QuoteDao) ByFriendlyName(op daoOperator.Operator, friendlyName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(%s.friendly_name) LIKE UPPER('%%%s%%')", d.TableName(), strings.TrimSpace(friendlyName)))
	}
}
func (d *QuoteDao) ByClientName(op daoOperator.Operator, name string) func(db *gorm.DB) *gorm.DB {
	clientTableName := clientDao.GetInstance().TableName()

	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(%s.name) LIKE UPPER('%%%s%%')", clientTableName, strings.TrimSpace(name)))
	}
}
func (d *QuoteDao) ByClientSurname(op daoOperator.Operator, surname string) func(db *gorm.DB) *gorm.DB {
	clientTableName := clientDao.GetInstance().TableName()

	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(%s.surname) LIKE UPPER('%%%s%%')", clientTableName, strings.TrimSpace(surname)))
	}
}
func (d *QuoteDao) ByLocality(op daoOperator.Operator, locality string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(%s.locality) LIKE UPPER('%%%s%%')", d.TableName(), strings.TrimSpace(locality)))
	}
}
func (d *QuoteDao) ByCountry(op daoOperator.Operator, country string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(%s.country) LIKE UPPER('%%%s%%')", d.TableName(), strings.TrimSpace(country)))
	}
}

//#endregion

//#endregion
