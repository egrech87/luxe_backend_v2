package colorModel

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/colorEntity"
	"Quote/model/baseModel"
	"context"
)

type ColorModel struct {
	*colorEntity.ColorEntity
}

func New() *ColorModel {
	return new(ColorModel)
}

func (m *ColorModel) GetEntity() baseEntity.BaseEntityI {
	return m.ColorEntity
}
func (m *ColorModel) SetEntity(value baseEntity.BaseEntityI) {
	m.ColorEntity = value.(*colorEntity.ColorEntity)
}

func (m *ColorModel) ToJson() map[string]any {
	return m.ColorEntity.ToJson()
}

func (m *ColorModel) FillFromJson(json map[string]any) error {
	if m.ColorEntity == nil {
		m.ColorEntity = colorEntity.New()
	}

	return m.ColorEntity.FillFromJson(json)
}

func (m *ColorModel) Clone() (baseModel.BaseModelI, error) {
	clone := New()
	clonedEntity, err := m.ColorEntity.Clone()
	if err != nil {
		return nil, err
	}

	clone.ColorEntity = clonedEntity.(*colorEntity.ColorEntity)
	return clone, nil
}

func (m *ColorModel) Validate(ctx context.Context) error {
	if err := m.GetEntity().Validate(); err != nil {
		return err
	}
	return nil
}

func (m *ColorModel) GetUpdates(oldModel baseModel.BaseModelI) map[string]any {
	return m.GetEntity().GetUpdates(oldModel.GetEntity())
}
