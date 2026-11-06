package glassModel

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/glassEntity"
	"Quote/model/baseModel"
	"context"
)

type GlassModel struct {
	*glassEntity.GlassEntity
}

func New() *GlassModel {
	return new(GlassModel)
}

func (m *GlassModel) GetEntity() baseEntity.BaseEntityI {
	return m.GlassEntity
}
func (m *GlassModel) SetEntity(value baseEntity.BaseEntityI) {
	m.GlassEntity = value.(*glassEntity.GlassEntity)
}

func (m *GlassModel) ToJson() map[string]any {
	return m.GlassEntity.ToJson()
}

func (m *GlassModel) FillFromJson(json map[string]any) error {
	if m.GlassEntity == nil {
		m.GlassEntity = glassEntity.New()
	}

	return m.GlassEntity.FillFromJson(json)
}

func (m *GlassModel) Clone() (baseModel.BaseModelI, error) {
	clone := New()
	clonedEntity, err := m.GlassEntity.Clone()
	if err != nil {
		return nil, err
	}

	clone.GlassEntity = clonedEntity.(*glassEntity.GlassEntity)
	return clone, nil
}

func (m *GlassModel) Validate(ctx context.Context) error {
	if err := m.GetEntity().Validate(); err != nil {
		return err
	}

	return nil
}

func (m *GlassModel) GetUpdates(oldModel baseModel.BaseModelI) map[string]any {
	return m.GetEntity().GetUpdates(oldModel.GetEntity())
}
