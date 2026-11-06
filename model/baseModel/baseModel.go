package baseModel

import (
	"Quote/db/entity/baseEntity"
	"context"
)

type BaseModelI interface {
	GetEntity() baseEntity.BaseEntityI
	SetEntity(baseEntity.BaseEntityI)
	ToJson() map[string]any
	FillFromJson(map[string]any) error
	Clone() (BaseModelI, error)
	Validate(context.Context) error
	GetUpdates(BaseModelI) map[string]any
}
