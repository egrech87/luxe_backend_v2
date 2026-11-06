package quoteModel

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/quoteEntity"
	"Quote/model/baseModel"
	"context"
)

type QuoteModel struct {
	*quoteEntity.QuoteEntity
}

func New() *QuoteModel {
	return new(QuoteModel)
}

func (m *QuoteModel) GetEntity() baseEntity.BaseEntityI {
	return m.QuoteEntity
}
func (m *QuoteModel) SetEntity(value baseEntity.BaseEntityI) {
	m.QuoteEntity = value.(*quoteEntity.QuoteEntity)
}

func (m *QuoteModel) ToJson() map[string]any {
	return m.QuoteEntity.ToJson()
}

func (m *QuoteModel) FillFromJson(json map[string]any) error {
	if m.QuoteEntity == nil {
		m.QuoteEntity = quoteEntity.New()
	}

	return m.QuoteEntity.FillFromJson(json)
}

func (m *QuoteModel) Clone() (baseModel.BaseModelI, error) {
	clone := New()
	clonedEntity, err := m.QuoteEntity.Clone()
	if err != nil {
		return nil, err
	}

	clone.QuoteEntity = clonedEntity.(*quoteEntity.QuoteEntity)
	return clone, nil
}

func (m *QuoteModel) Validate(ctx context.Context) error {
	if err := m.GetEntity().Validate(); err != nil {
		return err
	}
	return nil
}

func (m *QuoteModel) GetUpdates(oldModel baseModel.BaseModelI) map[string]any {
	return m.GetEntity().GetUpdates(oldModel.GetEntity())
}
