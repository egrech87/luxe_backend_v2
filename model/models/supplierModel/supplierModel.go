package supplierModel

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/supplierEntity"
	"Quote/model/baseModel"
	"context"
)

type SupplierModel struct {
	*supplierEntity.SupplierEntity
}

func New() *SupplierModel {
	return new(SupplierModel)
}

func (m *SupplierModel) GetEntity() baseEntity.BaseEntityI {
	return m.SupplierEntity
}
func (m *SupplierModel) SetEntity(value baseEntity.BaseEntityI) {
	m.SupplierEntity = value.(*supplierEntity.SupplierEntity)
}

func (m *SupplierModel) ToJson() map[string]any {
	return m.SupplierEntity.ToJson()
}

func (m *SupplierModel) FillFromJson(json map[string]any) error {
	if m.SupplierEntity == nil {
		m.SupplierEntity = supplierEntity.New()
	}

	return m.SupplierEntity.FillFromJson(json)
}

func (m *SupplierModel) Clone() (baseModel.BaseModelI, error) {
	clone := New()
	clonedEntity, err := m.SupplierEntity.Clone()
	if err != nil {
		return nil, err
	}

	clone.SupplierEntity = clonedEntity.(*supplierEntity.SupplierEntity)
	return clone, nil
}

func (m *SupplierModel) Validate(ctx context.Context) error {
	if err := m.GetEntity().Validate(); err != nil {
		return err
	}
	return nil
}

func (m *SupplierModel) GetUpdates(oldModel baseModel.BaseModelI) map[string]any {
	return m.GetEntity().GetUpdates(oldModel.GetEntity())
}
