package clientModel

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/clientEntity"
	"Quote/model/baseModel"
	"context"
)

type ClientModel struct {
	*clientEntity.ClientEntity
}

func New() *ClientModel {
	return new(ClientModel)
}

func (m *ClientModel) GetEntity() baseEntity.BaseEntityI {
	return m.ClientEntity
}
func (m *ClientModel) SetEntity(value baseEntity.BaseEntityI) {
	m.ClientEntity = value.(*clientEntity.ClientEntity)
}

func (m *ClientModel) ToJson() map[string]any {
	return m.ClientEntity.ToJson()
}

func (m *ClientModel) FillFromJson(json map[string]any) error {
	if m.ClientEntity == nil {
		m.ClientEntity = clientEntity.New()
	}

	return m.ClientEntity.FillFromJson(json)
}

func (m *ClientModel) Clone() (baseModel.BaseModelI, error) {
	clone := New()
	clonedEntity, err := m.ClientEntity.Clone()
	if err != nil {
		return nil, err
	}

	clone.ClientEntity = clonedEntity.(*clientEntity.ClientEntity)
	return clone, nil
}

func (m *ClientModel) Validate(ctx context.Context) error {
	if err := m.GetEntity().Validate(); err != nil {
		return err
	}

	return nil
}

func (m *ClientModel) GetUpdates(oldModel baseModel.BaseModelI) map[string]any {
	return m.GetEntity().GetUpdates(oldModel.GetEntity())
}
