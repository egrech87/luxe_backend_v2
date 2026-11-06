package apertureModel

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/apertureEntity"
	"Quote/model/baseModel"
	"context"
)

type ApertureModel struct {
	*apertureEntity.ApertureEntity
}

func New() *ApertureModel {
	return new(ApertureModel)
}

func (m *ApertureModel) GetEntity() baseEntity.BaseEntityI {
	return m.ApertureEntity
}

func (m *ApertureModel) SetEntity(value baseEntity.BaseEntityI) {
	m.ApertureEntity = value.(*apertureEntity.ApertureEntity)
}

func (m *ApertureModel) ToJson() map[string]any {
	return m.ApertureEntity.ToJson()
}

func (m *ApertureModel) FillFromJson(json map[string]any) error {
	if m.ApertureEntity == nil {
		m.ApertureEntity = apertureEntity.New()
	}

	return m.ApertureEntity.FillFromJson(json)
}

func (m *ApertureModel) Clone() (baseModel.BaseModelI, error) {
	clone := New()
	clonedEntity, err := m.ApertureEntity.Clone()
	if err != nil {
		return nil, err
	}

	clone.ApertureEntity = clonedEntity.(*apertureEntity.ApertureEntity)
	return clone, nil
}

func (m *ApertureModel) Validate(ctx context.Context) error {
	if err := m.GetEntity().Validate(); err != nil {
		return err
	}
	return nil
}

func (m *ApertureModel) GetUpdates(oldModel baseModel.BaseModelI) map[string]any {
	return m.GetEntity().GetUpdates(oldModel.GetEntity())
}
