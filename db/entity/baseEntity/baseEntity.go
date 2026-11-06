package baseEntity

import (
	"Quote/utils/mapUtils"
	"fmt"
	"time"

	"github.com/fatih/structs"
)

type BaseEntityI interface {
	GetID() uint64
	GetCreatedAt() time.Time
	GetCreatedBy() uint64
	GetUpdatedAt() time.Time
	GetUpdatedBy() uint64
	GetDeletedAt() *time.Time
	GetDeletedBy() uint64
	SetID(uint64)
	SetCreatedAt(time.Time)
	SetCreatedBy(uint64)
	SetUpdatedAt(time.Time)
	SetUpdatedBy(uint64)
	SetDeletedAt(*time.Time)
	SetDeletedBy(uint64)
	IsDeleted() bool

	TableName() string
	FieldNames() []string
	ToJson() map[string]any
	FillFromJson(json map[string]any) error
	FillFromDb(dbMap map[string]any) error
	Clone() (BaseEntityI, error)
	Validate() error
	GetUpdates(oldItem BaseEntityI) map[string]any
}

type BaseEntity struct {
	ID        uint64     `json:"id" gorm:"column:id; primaryKey; not null; auto_increment" structs:"id"`
	CreatedAt time.Time  `json:"created_at" structs:"created_at"`
	CreatedBy uint64     `json:"created_by" gorm:"column:created_by; not null" structs:"created_by"`
	UpdatedAt time.Time  `json:"updated_at" structs:"updated_at"`
	UpdatedBy uint64     `json:"updated_by" gorm:"column:updated_by; not null" structs:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" structs:"deleted_at"`
	DeletedBy uint64     `json:"deleted_by" gorm:"column:deleted_by" structs:"deleted_by"`
}

func (b *BaseEntity) GetID() uint64 {
	return b.ID
}
func (b *BaseEntity) GetCreatedAt() time.Time {
	return b.CreatedAt
}
func (b *BaseEntity) GetCreatedBy() uint64 {
	return b.CreatedBy
}
func (b *BaseEntity) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}
func (b *BaseEntity) GetUpdatedBy() uint64 {
	return b.UpdatedBy
}
func (b *BaseEntity) GetDeletedAt() *time.Time {
	return b.DeletedAt
}
func (b *BaseEntity) GetDeletedBy() uint64 {
	return b.DeletedBy
}
func (b *BaseEntity) SetID(value uint64) {
	b.ID = value
}
func (b *BaseEntity) SetCreatedAt(value time.Time) {
	b.CreatedAt = value
}
func (b *BaseEntity) SetCreatedBy(value uint64) {
	b.CreatedBy = value
}
func (b *BaseEntity) SetUpdatedAt(value time.Time) {
	b.UpdatedAt = value
}
func (b *BaseEntity) SetUpdatedBy(value uint64) {
	b.UpdatedBy = value
}
func (b *BaseEntity) SetDeletedAt(value *time.Time) {
	b.DeletedAt = value
}
func (b *BaseEntity) SetDeletedBy(value uint64) {
	b.DeletedBy = value
}

func (b *BaseEntity) IsDeleted() bool {
	return b.GetDeletedAt() != nil
}

func (*BaseEntity) TableName() string {
	return "base"
}
func (*BaseEntity) FieldNames() []string {
	return []string{"id", "created_at", "created_by", "updated_at", "updated_by", "deleted_at", "deleted_by"}
}

func (b *BaseEntity) ToJson() map[string]any {
	return structs.Map(b)
}

func (b *BaseEntity) FillFromJson(json map[string]any) error {
	_, err := mapUtils.FieldAnyToType(json, "id", fmt.Sprintf("%T", b.ID))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "created_by", fmt.Sprintf("%T", b.CreatedBy))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "created_at", fmt.Sprintf("%T", b.CreatedAt))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "updated_by", fmt.Sprintf("%T", b.UpdatedBy))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "updated_at", fmt.Sprintf("%T", b.UpdatedAt))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "deleted_by", fmt.Sprintf("%T", b.DeletedBy))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "deleted_at", fmt.Sprintf("%T", b.DeletedAt))
	if err != nil {
		return err
	}

	return mapUtils.ToStruct_usingJsonTags(json, b)
}
